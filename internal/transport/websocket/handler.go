package websocket

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const (
	timeout = 60
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
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			d, err := time.ParseDuration(fmt.Sprint(timeout/2) + "s")
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			time.Sleep(d)
		}
	}
}
