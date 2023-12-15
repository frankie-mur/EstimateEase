package server

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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

		msgType, payload, err := s.conn.ReadMessage()

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

		fmt.Printf("msgType: %v payload: %v\n", msgType, string(payload))
		fmt.Printf("Event: %v\n", event)
		//Update the user vote map for that specific user
		//NOTE: this is concurrent safe
		room.VoteMap.Update(s.name, event.Payload)

		htmlResponse := buildHTMLResponse(room)

		go s.Publisher.Broadcast(htmlResponse)
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
			log.Printf("message sent")

		//Wait on ticker to avoid subscriber from timing out
		case <-ticker.C:
			log.Printf("ping")
			//Send ping to client
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

// This function is used to build the HTML response for the room page
// Updating the vote map
func buildHTMLResponse(room *Room) string {
	sortedNames := room.VoteMap.SortedNames()

	trData := ""
	for _, name := range sortedNames {
		trData += fmt.Sprintf("<tr><td> %v </td><td> %v </td></tr>", name, room.VoteMap.VoteMap[name])
	}

	return fmt.Sprintf(`
	<div id="room-data">
	<div class="overflow-x-auto">
     <table class="table table-zebra">
     <!-- head -->
     <thead>
      <tr>
        <th>Name</th>
        <th>Vote</th>
      </tr>
     </thead>
      <tbody>
       %v
      </tbody>
     </table>
   </div> 
   <div>
   `, trData)
}
