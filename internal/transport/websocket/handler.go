package websocket

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func getNotificationsHandler(serv services.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params []services.Param
		params = append(params, services.Filter("notify_id", c.Params("notifyId"), true))
		limit, err := strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		params = append(params, services.Limit(limit))

		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		params = append(params, services.Limit(offset))

		status := c.Query("status")
		if status != "" {
			params = append(params, services.Filter("status", status, true))
		}

		notifications, err := serv.Get(context.Background(), params...)
		fmt.Printf("%+v", notifications)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
