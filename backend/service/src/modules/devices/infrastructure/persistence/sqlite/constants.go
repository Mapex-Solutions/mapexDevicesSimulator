package sqlite

// tableDevices is the SQLite table backing the devices module.
const tableDevices = "devices"

// DDL creates the devices table. Scalar columns are queryable; config,
// attributes and events hold serialized JSON. Idempotent: runs on every boot.
const DDL = `CREATE TABLE IF NOT EXISTS devices (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	device_id TEXT NOT NULL,
	protocol_id TEXT NOT NULL,
	enabled INTEGER NOT NULL,
	store_logs INTEGER NOT NULL,
	config TEXT,
	attributes TEXT,
	events TEXT,
	created TEXT NOT NULL
)`
