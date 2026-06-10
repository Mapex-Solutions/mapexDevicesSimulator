package sqlite

import (
	"context"
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

// New builds the log repository over the sqlite model bound to the logs table.
func New(mgr *sqliteManager.SQLiteManager) repositories.LogRepository {
	return &adapter{model: sqliteModel.New[entities.Log](mgr.DB(), tableLogs, sqliteModel.Config{})}
}

// Insert persists one message; the model fills id and created.
func (a *adapter) Insert(ctx context.Context, l *entities.Log) (*entities.Log, error) {
	return a.model.CreateOne(ctx, l)
}

// ListPage returns the matching logs newest-first plus the total matching count.
func (a *adapter) ListPage(ctx context.Context, f repositories.LogFilter) ([]entities.Log, int, error) {
	where, args := buildWhere(f)
	db := a.model.DIRECT()

	var total int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+tableLogs+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	query := "SELECT id, protocol, device_id, device_name, direction, kind, summary, status, payload, created FROM " +
		tableLogs + where + " ORDER BY created DESC LIMIT ? OFFSET ?"
	rows, err := db.QueryContext(ctx, query, append(args, limit, f.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]entities.Log, 0, limit)
	for rows.Next() {
		l, err := scanLog(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, l)
	}
	return items, total, rows.Err()
}

// buildWhere turns the filter into an equality + free-text WHERE clause. Column
// names are fixed literals (not user input), so only bound values reach SQL.
func buildWhere(f repositories.LogFilter) (string, []any) {
	conds := make([]string, 0, 5)
	args := make([]any, 0, 7)
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
	if f.Q != "" {
		conds = append(conds, "(summary LIKE ? OR payload LIKE ? OR device_name LIKE ?)")
		like := "%" + f.Q + "%"
		args = append(args, like, like, like)
	}
	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

// rowScanner is satisfied by *sql.Rows, kept narrow for scanLog.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanLog reads one row into a Log, parsing the stored created text.
func scanLog(rows rowScanner) (entities.Log, error) {
	var (
		l       entities.Log
		status  *string
		payload *string
		created string
	)
	if err := rows.Scan(&l.ID, &l.Protocol, &l.DeviceID, &l.DeviceName, &l.Direction, &l.Kind, &l.Summary, &status, &payload, &created); err != nil {
		return entities.Log{}, err
	}
	if status != nil {
		l.Status = *status
	}
	if payload != nil {
		l.Payload = *payload
	}
	if created != "" {
		if t, err := time.Parse(timeLayout, created); err == nil {
			l.Created = t.UTC()
		}
	}
	return l, nil
}
