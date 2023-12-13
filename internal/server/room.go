package server

import "github.com/google/uuid"

// A vote will represent
// Each users vote of a room
// Mapping username to their vote
type VoteMap map[string]string

// Room represents a room in the server (or session)
// Each room will have one publisher
// a publisher can have one or more subscribers
type Room struct {
	Id       uuid.UUID  `json:"id"`
	RoomName string     `json:"roomName"`
	Pub      *Publisher `json:"-"`
	VoteMap  VoteMap    `json:"_"`
}

type RoomsList map[*Room]bool

func NewRoom(name string) *Room {
	//TODO: How do we want to handle the id?
	return &Room{
		Id:       uuid.New(),
		RoomName: name,
		Pub:      NewPublisher(),
		VoteMap:  make(VoteMap),
	}
}
