package websocket

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Alieksieiev0/notification-service/internal/models"
	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const (
	timeout = 60
)

func getNotificationsHandler(serv services.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := []services.Param{services.Filter("notify_id", c.Params("notifyId"), true)}

		if l := c.Query("limit"); l != "" {
			convL, err := strconv.Atoi(l)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{"error": err.Error()})
			}
			params = append(params, services.Limit(convL))
		}

		if o := c.Query("offset"); o != "" {
			convO, err := strconv.Atoi(o)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{"error": err.Error()})
			}
			params = append(params, services.Limit(convO))
		}

		if ca := c.Query("created_at"); ca != "" {
			params = append(params, services.GTE("created_at", ca))
		}

		params = append(
			params,
			services.Order(c.Query("sort_by", "Id"), c.Query("order_by", "asc")),
		)

		notifications, err := serv.Get(context.Background(), params...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(notifications)
	}
}

func reviewHandler(serv services.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		fmt.Println("----------")
		id := c.Params("id")
		user, err := serv.GetById(
			c.Context(),
			id,
			services.Filter("status", models.NewNotificationStatus, true),
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		user.Status = models.ReviewedNoificationStatus

		err = serv.Save(context.Background(), user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		fmt.Println("----------")
		c.Status(fiber.StatusOK)
		return nil
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
