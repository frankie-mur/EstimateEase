package server

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/frankie-mur/EstimateEase/ui/components"
	"github.com/teris-io/shortid"
)

// Room represents a room in the server (or session)
// Each room will have one Publisher
// a Publisher can have one or more subscribers
// A Room hold the VoteMap used to store all subscribers names maped to their vote
type Room struct {
	Id               string     `json:"id"`
	RoomName         string     `json:"roomName"`
	Pub              *Publisher `json:"-"`
	VoteMap          *Votes     `json:"_"`
	VotesReveledFlag bool
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
	for room := range r.Rooms {
		if room.Id == id {
			return room, true
		}
	}
	return nil, false
}

// OnSubRemoved will update vote map for all users
func (r *Room) OnSubRemoved(subscriber *Subscriber) {
	r.VoteMap.Remove(subscriber.name)
	//Render voteMap component and broadcast to all subscribers
	voteMap := components.VoteMapData{
		SortedNames: r.VoteMap.SortNames(),
		VoteMap:     r.VoteMap.VoteMap,
		ShowVotes:   r.VotesReveledFlag,
	}
	data, err := RenderComponentToString(components.VotingGrid(voteMap), context.TODO())
	if err != nil {
		fmt.Print("Error rendering component: ", err)
		return
	}
	go r.Pub.Broadcast(data)
}

// Calculate the average all votes in a room
func (r *Room) calculateRoomStats() string {
	sum := 0
	for _, vote := range r.VoteMap.VoteMap {
		v, _ := strconv.Atoi(vote)
		sum += v
	}

	// Convert sum to float64 for floating-point division
	avg := float64(sum) / float64(len(r.VoteMap.VoteMap))

	res := fmt.Sprintf("%.2f", avg)

	return res
}

func (r *Room) Size() string {
	return fmt.Sprintf("%v", len(r.Pub.Subs))
}
