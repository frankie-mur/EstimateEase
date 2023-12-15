package server

import (
	"estimate-ease/internal/data"

	"github.com/google/uuid"
)

// Room represents a room in the server (or session)
// Each room will have one Publisher
// a Publisher can have one or more subscribers
// A Room hold the VoteMap used to store all subscribers names maped to their vote
type Room struct {
	//TODO: move away from UUID as it is too long and not very human readable/good for url's
	Id       uuid.UUID   `json:"id"`
	RoomName string      `json:"roomName"`
	Pub      *Publisher  `json:"-"`
	VoteMap  *data.Votes `json:"_"`
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
		VoteMap:  data.NewVoteMap(),
	}
}

func (r *RoomsList) Is(id uuid.UUID) (*Room, bool) {
	//Need to dereference the roomList to iterate over
	for room := range *r {
		if room.Id == id {
			return room, true
		}
	}
	return nil, false
}
