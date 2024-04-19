package app

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/notification-service/internal/config"
	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/Alieksieiev0/notification-service/internal/transport/kafka"
	"github.com/Alieksieiev0/notification-service/internal/transport/websocket"
	"github.com/joho/godotenv"
)

func Run() {
	var (
		websocketServerAddr = flag.String(
			"websocket-server",
			":3003",
			"listen address of websocket server",
		)
		kafkaAddr = flag.String("kafka", "kafka:9094", "address of kafka")
	)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.Database()
	if err != nil {
		log.Fatal(err)
	}

	_ = websocketServerAddr
	c := kafka.NewConsumer([]string{*kafkaAddr})
	err = c.Consume(services.NewService(db), websocket.NewTransferer())
	if err != nil {
		log.Fatal(err)
	}

}
