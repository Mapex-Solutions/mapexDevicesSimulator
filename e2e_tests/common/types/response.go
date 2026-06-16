// Package types holds the wire shapes the e2e decodes from the simulator.
package types

import "encoding/json"

// Envelope is the simulator's standard response wrapper: every REST reply is
// `{ status, errors, data }`. Data is kept raw so each caller decodes it into
// the concrete payload (a Device, a LogPage, ...) it expects.
type Envelope struct {
	Status int             `json:"status"`
	Errors []string        `json:"errors"`
	Data   json.RawMessage `json:"data"`
}
