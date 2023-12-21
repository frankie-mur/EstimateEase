package server

import (
	"fmt"
	"sync"

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

type RoomList struct {
	Rooms map[*Room]bool
	sync.RWMutex
}

func NewRoom(name string) *Room {
	return &Room{
		Id:       generateId(),
		RoomName: name,
		Pub:      NewPublisher(),
		VoteMap:  NewVoteMap(),
	}
}

func NewRoomList() *RoomList {
	return &RoomList{
		Rooms: make(map[*Room]bool),
	}
}

func generateId() string {
	// generate a short unique id
	return shortid.MustGenerate()
}

func (r *RoomList) Is(id string) (*Room, bool) {
	//Need to dereference the roomList to iterate over
	fmt.Printf("Looking for room with id %s", id)
	for room := range r.Rooms {
		fmt.Printf("Room % v", room)
		if room.Id == id {
			return room, true
		}
	}
	return nil, false
}
