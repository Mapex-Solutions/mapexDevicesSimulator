package ports

import (
	"simulator/service/src/modules/console/application/dtos"
)

// Publisher broadcasts a console frame to every connected WebSocket client. The
// simulation engine depends on this port to stream live activity without
// importing the WebSocket transport.
type Publisher interface {
	Publish(msg dtos.ConsoleMessage)
}
