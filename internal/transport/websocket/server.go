package websocket

import (
	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketServer struct {
	app  *fiber.App
	addr string
}

func NewServer(app *fiber.App, addr string) *WebsocketServer {
	return &WebsocketServer{
		app:  app,
		addr: addr,
	}
}

func (ws *WebsocketServer) Start(serv services.Service, trans Transferer) error {
	go trans.Run()
	ws.app.Use("/listen", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			// test with false
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.app.Get("/notifications/:notifyId", getNotificationsHandler(serv))
	ws.app.Get("/listen/:notifyId", websocket.New(listenHandler(trans)))

	return ws.app.Listen(ws.addr)
}
