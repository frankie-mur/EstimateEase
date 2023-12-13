package server

import "github.com/google/uuid"

// Room represents a room in the server (or session)
// Each room will have one publisher
// a publisher can have one or more subscribers
type Room struct {
	Id       uuid.UUID  `json:"id"`
	RoomName string     `json:"roomName"`
	Pub      *Publisher `json:"-"`
}

type RoomsList map[*Room]bool

func NewRoom(name string) *Room {
	//TODO: How do we want to handle the id?
	return &Room{
		Id:       uuid.New(),
		RoomName: name,
		Pub:      NewPublisher(),
	}
}
