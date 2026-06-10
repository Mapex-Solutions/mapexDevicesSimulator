package gateways

import (
	"encoding/json"

	timeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/time"
)

// Gateway is the wire shape of a LoRaWAN gateway (the Basics Station / UDP link to
// the LNS). `link` is passed through as raw JSON (its used fields depend on the
// link protocol).
type Gateway struct {
	ID          string             `json:"id"`
	Created     *timeUtil.NullTime `json:"created"`
	Name        string             `json:"name"`
	EUI         string             `json:"eui"`
	Enabled     bool               `json:"enabled"`
	Region      string             `json:"region"`
	Description string             `json:"description"`
	Link        json.RawMessage    `json:"link"`
}

// SetCreated lets mapper.EntityToDto populate the creation timestamp.
func (g *Gateway) SetCreated(t *timeUtil.NullTime) { g.Created = t }

// GatewayIDParam binds the :id path parameter on the gateway by-id routes.
type GatewayIDParam struct {
	ID string `params:"id" validate:"required"`
}

// GatewayInput is the create/update body.
type GatewayInput struct {
	Name        string          `json:"name" validate:"required"`
	EUI         string          `json:"eui" validate:"required"`
	Enabled     bool            `json:"enabled"`
	Region      string          `json:"region" validate:"required,oneof=EU868 US915 AU915 AS923 CN470 IN865 KR920 RU864"`
	Description string          `json:"description"`
	Link        json.RawMessage `json:"link"`
}
