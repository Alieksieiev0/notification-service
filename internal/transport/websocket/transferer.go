package websocket

import (
	"fmt"

	"github.com/Alieksieiev0/notification-service/internal/models"
	"github.com/gofiber/contrib/websocket"
)

type Transferer interface {
	Run()
	Pass(notification *models.Notification)
	AddConnection(conn *websocket.Conn, notifyId string)
}

func NewTransferer() Transferer {
	return &transferer{
		notification: make(chan *models.Notification),
		connections:  make(map[string]*websocket.Conn),
	}
}

type transferer struct {
	notification chan *models.Notification
	connections  map[string]*websocket.Conn
}

func (w *transferer) Run() {
	for {
		msg := <-w.notification
		fmt.Printf("msg: %+v", msg)
		c, ok := w.connections[msg.NotifyId]
		fmt.Println(ok)
		if !ok {
			continue
		}

		if err := c.WriteJSON(msg); err != nil {
			fmt.Println(err)
			delete(w.connections, msg.NotifyId)
		}
	}
}

func (w *transferer) Pass(notification *models.Notification) {
	w.notification <- notification
}

func (w *transferer) AddConnection(conn *websocket.Conn, notifyId string) {
	w.connections[notifyId] = conn
}
