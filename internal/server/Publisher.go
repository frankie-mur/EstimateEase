package server

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Publisher struct {
	//List of all subscriptions to a producer
	Subs SubscriberList
	//Used to lock before editing subs
	//As maps are not concurrent safe
	sync.Mutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		Subs: make(SubscriberList),
	}
}

func (p *Publisher) AddSubscriber(sub *Subscriber) {
	p.Lock()
	defer p.Unlock()
	p.Subs[sub] = true
}

func (p *Publisher) RemoveSubscriber(sub *Subscriber) {
	p.Lock()
	defer p.Unlock()

	//Check if the subscriber is in the subscriber list
	if _, ok := p.Subs[sub]; ok {
		//If it is we close the subscriber connection
		sub.conn.Close()
		//And delete from subscriber list
		delete(p.Subs, sub)
	}
}

func (p *Publisher) Broadcast(voteMap *Votes) error {
	//For each user we want to send them a updated list of
	//All the current votes
	htmlResponse :=
		fmt.Sprintf("<div id=\"room-data\">Current Results: %v</div>", voteMap)

	for subs := range p.Subs {
		fmt.Printf("Braoadcasting to %v with data %v\n", subs.conn.RemoteAddr(), htmlResponse)
		subs.egress <- []byte(htmlResponse)
	}
	return nil
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
