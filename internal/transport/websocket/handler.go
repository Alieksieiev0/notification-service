package websocket

import (
	"context"
	"log"

	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func getNotificationsHandler(serv services.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notifyId := c.Params("notifyId")
		notifications, err := serv.GetByNotifyId(context.Background(), notifyId)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(notifications)
	}
}

func listenHandler(trans Transferer) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		notifyId := c.Params("notifyId")
		trans.AddConnection(c, notifyId)

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
