package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

// SubscriberList Used to help manage subscribers
type SubscriberList map[*Subscriber]bool

type Subscriber struct {
	conn      *websocket.Conn
	publisher *Publisher
	// egress is used to avoid concurrent writes on websocket conn
	egress chan []byte
}

func NewSubscriber(conn *websocket.Conn, publisher *Publisher) *Subscriber {
	return &Subscriber{
		conn:      conn,
		publisher: publisher,
		egress:    make(chan []byte),
	}
}

// Reads a websocket message setting read limits and pong handler
// for safety and durability
func (s *Subscriber) readMessage() {
	defer func() {
		// clean up connection
		s.publisher.removeSubscriber(s)
	}()
	if err := s.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		fmt.Println(err)
		return
	}
	//TODO: set this limit to htmx websocket request size
	s.conn.SetReadLimit(512)

	s.conn.SetPongHandler(s.pongHandler)

	for {
		msgType, payload, err := s.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		fmt.Printf("msgType: %v payload: %v\n", msgType, string(payload))
		//Broadcast the message to all subscribers
		go func() {
			err := s.publisher.broadcast(payload)
			if err != nil {
				log.Printf("failed to broadcast message: %v", err)
			}
		}()
	}
}

func (s *Subscriber) writeMessages() {
	defer func() {
		s.publisher.removeSubscriber(s)
	}()
	ticker := time.NewTicker(pingInterval)

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

		//Wait on ticker to avoid subscriber from timing out
		case <-ticker.C:
			log.Printf("ping")
			//Send ping to clinet
			if err := s.conn.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Printf("failed to send ping: %v", err)
				return
			}
		}
	}
}

func (s *Subscriber) pongHandler(pongMsg string) error {
	log.Println("pong")
	return s.conn.SetReadDeadline(time.Now().Add(pongWait))
}
