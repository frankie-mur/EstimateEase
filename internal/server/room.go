package server

import "github.com/google/uuid"

// Room represents a room in the server (or session)
// Each room will have one publisher
// a publisher can have one or more subscribers
// A Room hold the VoteMap used to store all subscribers names maped to their vote
type Room struct {
	//TODO: move away from UUID as it is too long and not very human readable/good for url's
	Id       uuid.UUID  `json:"id"`
	RoomName string     `json:"roomName"`
	Pub      *Publisher `json:"-"`
	VoteMap  *Votes     `json:"_"`
}

// RoomPageData data needed for the room page html template
type RoomPageData struct {
	Room        *Room
	DisplayName string
}

type RoomsList map[*Room]bool

func NewRoom(name string) *Room {
	return &Room{
		Id:       uuid.New(),
		RoomName: name,
		Pub:      NewPublisher(),
		VoteMap:  NewVoteMap(),
	}
}
