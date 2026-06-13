package sqlite

// tableLogs is the SQLite table backing the logs module.
const tableLogs = "logs"

const (
	// defaultLimit is the page size when the query omits limit.
	defaultLimit = 50
	// maxLimit caps the page size so a query cannot pull the whole table.
	maxLimit = 500
)

// Migrations creates the logs table and the keyset index on (created, id) used
// by cursor pagination. Idempotent: runs on every boot. Columns added after the
// first release are reconciled by EnsureSchema, since SQLite has no
// ADD COLUMN IF NOT EXISTS.
var Migrations = []string{
	`CREATE TABLE IF NOT EXISTS logs (
		id TEXT PRIMARY KEY,
		protocol TEXT NOT NULL,
		device_id TEXT NOT NULL,
		device_name TEXT NOT NULL,
		event_name TEXT,
		direction TEXT NOT NULL,
		kind TEXT NOT NULL,
		summary TEXT NOT NULL,
		status TEXT,
		payload TEXT,
		response TEXT,
		created TEXT NOT NULL
	)`,
	`CREATE INDEX IF NOT EXISTS idx_logs_created_id ON logs(created, id)`,
}

// addedColumns are columns introduced after the first release. EnsureSchema adds
// any that an existing database is missing, so the table converges to the schema
// above without dropping history.
var addedColumns = []struct{ name, ddl string }{
	{"event_name", "event_name TEXT"},
	{"response", "response TEXT"},
}
