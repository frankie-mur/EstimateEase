package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/frankie-mur/EstimateEase/ui/components"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 1 * time.Second
	pingInterval = (pongWait * 9) / 10
)

// SubscriberList Used to help manage subscribers
type SubscriberList map[*Subscriber]bool

type Subscriber struct {
	conn      *websocket.Conn
	Publisher *Publisher
	// egress is used to avoid concurrent writes on websocket conn
	egress chan []byte
	name   string
}

type SubRemovedCallback interface {
	OnSubRemoved(subscriber *Subscriber)
}

func NewSubscriber(conn *websocket.Conn, room *Room, displayName string) *Subscriber {
	sub := &Subscriber{
		conn:      conn,
		Publisher: room.Pub,
		egress:    make(chan []byte),
		name:      displayName,
	}

	//Start subscriber background processes to
	//read and write messages
	//room is passed into ReadMessage as needed to broadcast messages to all subscribers
	go sub.ReadMessage(room)
	go sub.WriteMessages()

	return sub
}

// ReadMessage Reads a websocket message setting read limits and pong handler
//  1. SetReadDeadline to pongWait from now + pongWait
//  2. SetReadLimit to 512 bytes (max message size) to avoid large messages
//  3. SetPongHandler to pongHandler to handle pong messages
//  4. Read messages from the websocket connection
//  5. If the message of type Event then update the voteMap to the Event payload
//     and broadcast updated voteMap to all subscribers
//  6. If the message is a close message or error then close the connection
func (s *Subscriber) ReadMessage(room *Room) {
	defer func() {
		// clean up connection
		s.Publisher.RemoveSubscriber(s)
	}()
	if err := s.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		fmt.Println(err)
		return
	}
	s.conn.SetReadLimit(512)

	s.conn.SetPongHandler(s.pongHandler)

	for {
		var event Event

		_, payload, err := s.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		err = json.Unmarshal(payload, &event)
		if err != nil {
			fmt.Print("Error unmarshalling event: ", err)
			break
		}

		var avgOfVotes string
		switch event.Payload {
		case "show-votes":
			avgOfVotes = room.calculateRoomStats()
			room.VotesReveledFlag = true
		default:
			//Default payload is that of pressing a vote button
			room.VoteMap.Update(s.name, event.Payload)
		}

		//Render voteMap component and broadcast to all subscribers
		voteMap := components.VoteMapData{
			SortedNames: room.VoteMap.SortNames(),
			VoteMap:     room.VoteMap.VoteMap,
			ShowVotes:   room.VotesReveledFlag,
			Stats: []components.StatValues{{
				StatName:  "Total Users",
				StatValue: room.Size(),
			}, {
				StatName:  "Average",
				StatValue: avgOfVotes,
			}},
		}

		room.VotesReveledFlag = false

		data, err := RenderComponentToString(components.VotingGrid(voteMap), context.TODO())
		if err != nil {
			fmt.Print("Error rendering component: ", err)
			break
		}
		go s.Publisher.Broadcast(data)
	}
}

// WriteMessages writes messages to the websocket connection
// 1. Sets a ticker interval to pong the client
// 2. Reads messages from the egress channel
// 3. Sends messages to the client
func (s *Subscriber) WriteMessages() {
	defer func() {
		s.Publisher.RemoveSubscriber(s)
	}()
	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case msg, ok := <-s.egress:
			if !ok {
				if err := s.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed %v", err)
				}
				return
			}

			if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("failed to send message: %v", err)
			}

		//Wait on ticker to avoid subscriber from timing out
		case <-ticker.C:
			//Send ping to client
			if err := s.conn.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Printf("failed to send ping: %v", err)
				return
			}
		}
	}
}

func (s *Subscriber) pongHandler(pongMsg string) error {
	return s.conn.SetReadDeadline(time.Now().Add(pongWait))
}
