package persistence

// Schema lists the DDL statements applied at boot, in order. Each statement must
// be idempotent (CREATE TABLE IF NOT EXISTS ...) because it runs on every start.
// Modules append their tables here as they land.
var Schema = []string{}
