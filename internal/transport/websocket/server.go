package websocket

import (
	"net/http"

	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketServer struct {
	app  *fiber.App
	addr string
}

func (ws *WebsocketServer) Start(serv services.Service, trans Transferer) error {
	ws.app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// test with false
			c.Locals("allowed", true)
			if err := c.Next(); err != nil {
				return err
			}
		}
		return fiber.ErrUpgradeRequired
	})
	ws.app.Get("/notifications/:notifyId", websocket.New(getNotificationsHandler(serv, trans)))

	return http.ListenAndServe(ws.addr, nil)
}
