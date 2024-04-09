package websocket

import (
	"net/http"
)

const (
	subscriptionTopic = "subscriptions"
	postTopic         = "posts"
)

func Start(addr string, kafkaAddr string) error {
	http.HandleFunc("/notifications/posts", sendNotifications(postTopic, kafkaAddr))
	http.HandleFunc("/notifications/subscriptions", sendNotifications(subscriptionTopic, kafkaAddr))

	return http.ListenAndServe(addr, nil)
}
