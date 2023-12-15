package server

import (
	"github.com/teris-io/shortid"
)

// Room represents a room in the server (or session)
// Each room will have one Publisher
// a Publisher can have one or more subscribers
// A Room hold the VoteMap used to store all subscribers names maped to their vote
type Room struct {
	//TODO: move away from UUID as it is too long and not very human readable/good for url's
	Id       string     `json:"id"`
	RoomName string     `json:"roomName"`
	Pub      *Publisher `json:"-"`
	VoteMap  *Votes     `json:"_"`
}

type RoomsList map[*Room]bool

func NewRoom(name string) *Room {
	return &Room{
		Id:       generateId(),
		RoomName: name,
		Pub:      NewPublisher(),
		VoteMap:  NewVoteMap(),
	}
}

func generateId() string {
	// generate a short unique id
	return shortid.MustGenerate()
}

func (r *RoomsList) Is(id string) (*Room, bool) {
	//Need to dereference the roomList to iterate over
	for room := range *r {
		if room.Id == id {
			return room, true
		}
	}
	return nil, false
}
