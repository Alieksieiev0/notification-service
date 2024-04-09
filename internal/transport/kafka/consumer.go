package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

func Consume(conn *websocket.Conn, topic string, addr ...string) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  addr,
		GroupID:  topic + "-consumer",
		Topic:    topic,
		MaxBytes: 10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		body := map[string]interface{}{}
		err = json.Unmarshal(m.Value, &body)
		if err != nil {
			break
		}

		body["ID"] = string(m.Key)

		err = conn.WriteJSON(body)
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	return r.Close()
}
