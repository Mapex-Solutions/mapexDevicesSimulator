package entities

import (
	"encoding/json"
	"time"
)

// Gateway is the persisted LoRaWAN gateway. Scalar columns are queryable; the
// link block is stored as JSON text.
type Gateway struct {
	ID          string          `db:"id,pk"`
	Name        string          `db:"name"`
	EUI         string          `db:"eui"`
	Enabled     bool            `db:"enabled"`
	Region      string          `db:"region"`
	Description string          `db:"description"`
	Link        json.RawMessage `db:"link,json"`
	Created     time.Time       `db:"created"`
}

// GetCreated exposes the creation time for the entity-to-DTO mapper.
func (g *Gateway) GetCreated() time.Time { return g.Created }
