package sqlite

// tableGateways is the SQLite table backing the gateways module.
const tableGateways = "gateways"

// DDL creates the gateways table. Idempotent: runs on every boot.
const DDL = `CREATE TABLE IF NOT EXISTS gateways (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	eui TEXT NOT NULL,
	enabled INTEGER NOT NULL,
	region TEXT NOT NULL,
	description TEXT,
	link TEXT,
	created TEXT NOT NULL
)`
