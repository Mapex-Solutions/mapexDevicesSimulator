package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"

	"simulator/service/src/modules/console/infrastructure/hub"
)

// Register mounts the console WebSocket at /ws. A non-upgrade request is rejected;
// an upgraded connection subscribes to the hub and streams frames until it closes.
func Register(app *fiber.App, h *hub.Hub) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		// A plain HTTP hit on the websocket path flows through the global error
		// handler as the standard envelope, carrying the upgrade-required code.
		return &customErrors.ServerCustomError{
			Code:   fiber.StatusUpgradeRequired,
			Errors: []string{"websocket upgrade required"},
		}
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		cl := h.Register()
		defer h.Unregister(cl)

		// Read pump: discard inbound frames but detect the close so the write
		// pump can exit and the client is unregistered.
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					h.Unregister(cl)
					return
				}
			}
		}()

		// Write pump: forward broadcast frames until the channel closes.
		for data := range cl.Out() {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}))
}
