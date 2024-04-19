package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Alieksieiev0/notification-service/internal/models"
	"github.com/Alieksieiev0/notification-service/internal/services"
	"github.com/Alieksieiev0/notification-service/internal/transport/websocket"
	"github.com/segmentio/kafka-go"
)

const (
	postsTopic         = "posts"
	subscriptionsTopic = "subscriptions"
)

type Consumer interface {
	Consume(serv services.Service, trans websocket.Transferer) error
}

type consumer struct {
	addrs  []string
	topics []string
}

func NewConsumer(addrs []string, topics ...string) Consumer {
	if len(topics) == 0 {
		topics = append(topics, postsTopic, subscriptionsTopic)
	}

	return &consumer{
		addrs:  addrs,
		topics: topics,
	}
}

func (c *consumer) Consume(serv services.Service, trans websocket.Transferer) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.addrs,
		GroupID:     strings.Join(c.topics, "-") + "-consumer",
		GroupTopics: c.topics,
		MaxBytes:    10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		fmt.Println("-----waiting message------")
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("-----received message------")

		n := &models.Notification{}
		err = json.Unmarshal(m.Value, &n)
		if err != nil {
			break
		}
		n.NotifyId = string(m.Key)
		err = serv.Save(context.Background(), n)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		fmt.Println("-----before pass------")
		go trans.Pass(n)
		fmt.Println("-----afer pass------")
	}

	return r.Close()
}
