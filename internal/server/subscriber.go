package server

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Used to help manage client
type SubscriberList map[*Subscriber]bool

type Subscriber struct {
	conn      *websocket.Conn
	publisher *Publisher
	// egress is used to avoid concurrent writes on websocket conn
	egress chan []byte
}

func NewSubscriber(conn *websocket.Conn, publisher *Producer) *Subscriber {
	return &Subscriber{
		conn:      conn,
		publisher: publisher,
		egress:    make(chan []byte),
	}
}

// Start read message
func (s *Subscriber) readMessage() {
	defer func() {
		// clean up connection
		s.publisher.removeSubscriber(s)
	}()
	for {
		msgType, payload, err := s.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		fmt.Printf("msgType: %v payload: %v\n", msgType, string(payload))
	}
}

func (s *Subscriber) writeMessages() {
	defer func() {
		s.publisher.removeSubscriber(s)
	}()

	for {
		select {
		case message, ok := <-s.egress:
			if !ok {
				if err := s.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed %v", err)
				}
				return
			}

			if err := s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Printf("message sent")
		}
	}
}
