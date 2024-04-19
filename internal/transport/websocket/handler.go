package websocket

import (
	"context"
	"log"

	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func getNotificationsHandler(serv services.Service, trans Transferer) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		notifyId := c.Params("notifyId")
		notifications, err := serv.GetByNotifyId(context.Background(), notifyId)
		if err != nil {
			if err = c.WriteJSON(fiber.Map{"error": err.Error()}); err != nil {
				log.Println(err.Error())
			}
			return
		}

		if err = c.WriteJSON(notifications); err != nil {
			log.Println(err.Error())
			return
		}

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
