package websocket

import (
	"log"
	"net/http"

	"github.com/Alieksieiev0/notification-service/internal/transport/kafka"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func sendNotifications(topic string, addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		kafka.Consume(c, topic, addr)
	}
}
