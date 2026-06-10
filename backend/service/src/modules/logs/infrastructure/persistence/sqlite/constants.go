package sqlite

// tableLogs is the SQLite table backing the logs module.
const tableLogs = "logs"

const (
	// defaultLimit is the page size when the query omits limit.
	defaultLimit = 50
	// maxLimit caps the page size so a query cannot pull the whole table.
	maxLimit = 500
)

// Migrations creates the logs table and an index on created for the newest-first
// console history. Idempotent: runs on every boot.
var Migrations = []string{
	`CREATE TABLE IF NOT EXISTS logs (
		id TEXT PRIMARY KEY,
		protocol TEXT NOT NULL,
		device_id TEXT NOT NULL,
		device_name TEXT NOT NULL,
		direction TEXT NOT NULL,
		kind TEXT NOT NULL,
		summary TEXT NOT NULL,
		status TEXT,
		payload TEXT,
		created TEXT NOT NULL
	)`,
	`CREATE INDEX IF NOT EXISTS idx_logs_created ON logs(created)`,
}
