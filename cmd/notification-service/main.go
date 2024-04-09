package main

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/notification-service/internal/transport/websocket"
)

func main() {
	var (
		websocketServerAddr = flag.String(
			"websocket-server",
			":3003",
			"listen address of websocket server",
		)
		kafkaAddr = flag.String("kafka", "9092", "address of kafka")
	)

	err := websocket.Start(*websocketServerAddr, *kafkaAddr)
	if err != nil {
		log.Fatal(err)
	}
}
