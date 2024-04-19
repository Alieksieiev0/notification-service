package app

import (
	"flag"
	"fmt"
	"log"

	"github.com/Alieksieiev0/notification-service/internal/config"
	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/Alieksieiev0/notification-service/internal/transport/kafka"
	"github.com/Alieksieiev0/notification-service/internal/transport/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func Run() {
	var (
		websocketServerAddr = flag.String(
			"websocket-server",
			":3002",
			"listen address of websocket server",
		)
		kafkaAddr = flag.String("kafka", "kafka:9094", "address of kafka")
		app       = fiber.New()
		g         = new(errgroup.Group)
	)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.Database()
	if err != nil {
		log.Fatal(err)
	}

	serv := services.NewService(db)
	trans := websocket.NewTransferer()
	c := kafka.NewConsumer([]string{*kafkaAddr})
	g.Go(func() error {
		return c.Consume(serv, trans)
	})

	ws := websocket.NewServer(app, *websocketServerAddr)
	g.Go(func() error {
		return ws.Start(serv, trans)
	})

	if err := g.Wait(); err != nil {
		fmt.Println("----------------")
		log.Fatal(err)
	}
}
