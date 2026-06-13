package sqlite

import (
	"context"
	"database/sql"
	"encoding/base64"
	"strings"
	"time"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"
	sqliteModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/model"

	"simulator/service/src/modules/logs/domain/entities"
	"simulator/service/src/modules/logs/domain/repositories"
)

// Compile-time proof the adapter satisfies the repository port.
var _ repositories.LogRepository = (*adapter)(nil)

// timeLayout matches how the sqlite model stores time.Time, so created round-trips.
const timeLayout = "2006-01-02T15:04:05.999999999Z07:00"

// cursorSep separates the created and id parts inside an encoded cursor token.
const cursorSep = "\x00"

// New builds the log repository over the sqlite model bound to the logs table.
func New(mgr *sqliteManager.SQLiteManager) repositories.LogRepository {
	return &adapter{model: sqliteModel.New[entities.Log](mgr.DB(), tableLogs, sqliteModel.Config{})}
}

// EnsureSchema runs the base migrations and adds any columns an older database
// is missing. SQLite has no ADD COLUMN IF NOT EXISTS, so each add is guarded by
// a table_info check to keep the boot idempotent.
func EnsureSchema(ctx context.Context, mgr *sqliteManager.SQLiteManager) error {
	if err := mgr.Migrate(ctx, Migrations...); err != nil {
		return err
	}
	db := mgr.DB()
	existing, err := existingColumns(ctx, db)
	if err != nil {
		return err
	}
	for _, col := range addedColumns {
		if existing[col.name] {
			continue
		}
		if _, err := db.ExecContext(ctx, "ALTER TABLE "+tableLogs+" ADD COLUMN "+col.ddl); err != nil {
			return err
		}
	}
	return nil
}

// existingColumns returns the set of columns currently on the logs table.
func existingColumns(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, "PRAGMA table_info("+tableLogs+")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols := make(map[string]bool)
	for rows.Next() {
		var (
			cid, notnull, pk int
			name, ctype      string
			dflt             any
		)
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return nil, err
		}
		cols[name] = true
	}
	return cols, rows.Err()
}

// Insert persists one message; the model fills id and created.
func (a *adapter) Insert(ctx context.Context, l *entities.Log) (*entities.Log, error) {
	return a.model.CreateOne(ctx, l)
}

// ListPage returns the matching logs newest-first and the cursor for the next
// page. It fetches one row beyond the limit to learn whether a next page exists;
// when it does, the extra row is dropped and the limit-th row becomes the cursor.
func (a *adapter) ListPage(ctx context.Context, f repositories.LogFilter) ([]entities.Log, string, error) {
	where, args := buildWhere(f)
	limit := f.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	query := "SELECT id, protocol, device_id, device_name, event_name, direction, kind, summary, status, payload, response, created FROM " +
		tableLogs + where + " ORDER BY created DESC, id DESC LIMIT ?"
	rows, err := a.model.DIRECT().QueryContext(ctx, query, append(args, limit+1)...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	items := make([]entities.Log, 0, limit+1)
	for rows.Next() {
		l, err := scanLog(rows)
		if err != nil {
			return nil, "", err
		}
		items = append(items, l)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	next := ""
	if len(items) > limit {
		next = encodeCursor(items[limit-1])
		items = items[:limit]
	}
	return items, next, nil
}

// buildWhere turns the filter into the WHERE clause. Column names are fixed
// literals (never user input), so only bound values reach SQL. The cursor adds a
// keyset predicate matching the (created DESC, id DESC) order.
func buildWhere(f repositories.LogFilter) (string, []any) {
	conds := make([]string, 0, 8)
	args := make([]any, 0, 12)
	add := func(col, val string) {
		if val != "" {
			conds = append(conds, col+" = ?")
			args = append(args, val)
		}
	}
	add("protocol", f.Protocol)
	add("kind", f.Kind)
	add("direction", f.Direction)
	add("device_id", f.Device)
	if f.Event != "" {
		conds = append(conds, "event_name LIKE ?")
		args = append(args, "%"+f.Event+"%")
	}
	if f.DateFrom != "" {
		conds = append(conds, "created >= ?")
		args = append(args, f.DateFrom)
	}
	if f.DateTo != "" {
		conds = append(conds, "created <= ?")
		args = append(args, f.DateTo)
	}
	if f.Q != "" {
		conds = append(conds, "(summary LIKE ? OR payload LIKE ? OR device_name LIKE ?)")
		like := "%" + f.Q + "%"
		args = append(args, like, like, like)
	}
	if created, id, ok := decodeCursor(f.Cursor); ok {
		conds = append(conds, "(created < ? OR (created = ? AND id < ?))")
		args = append(args, created, created, id)
	}
	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

// encodeCursor builds the opaque keyset token for a row: its created text and id.
func encodeCursor(l entities.Log) string {
	raw := l.Created.Format(timeLayout) + cursorSep + l.ID
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// decodeCursor splits a token back into the created text and id; ok is false when
// the token is empty or malformed, so a bad cursor reads as "first page".
func decodeCursor(s string) (created, id string, ok bool) {
	if s == "" {
		return "", "", false
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(b), cursorSep, 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// rowScanner is satisfied by *sql.Rows, kept narrow for scanLog.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanLog reads one row into a Log, parsing the stored created text.
func scanLog(rows rowScanner) (entities.Log, error) {
	var (
		l         entities.Log
		eventName *string
		status    *string
		payload   *string
		response  *string
		created   string
	)
	if err := rows.Scan(&l.ID, &l.Protocol, &l.DeviceID, &l.DeviceName, &eventName, &l.Direction, &l.Kind, &l.Summary, &status, &payload, &response, &created); err != nil {
		return entities.Log{}, err
	}
	if eventName != nil {
		l.EventName = *eventName
	}
	if status != nil {
		l.Status = *status
	}
	if payload != nil {
		l.Payload = *payload
	}
	if response != nil {
		l.Response = *response
	}
	if created != "" {
		if t, err := time.Parse(timeLayout, created); err == nil {
			l.Created = t.UTC()
		}
	}
	return l, nil
}
