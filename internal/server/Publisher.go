package server

import (
	"fmt"
	"sync"
)

type Publisher struct {
	//List of all subscriptions to a producer
	subs SubscriberList
	//Used to lock before editing subs
	//As maps are not concurrent safe
	sync.Mutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		subs: make(SubscriberList),
	}
}

func (p *Publisher) addSubscriber(sub *Subscriber) {
	p.Lock()
	p.subs[sub] = true
	p.Unlock()
}

func (p *Publisher) removeSubscriber(sub *Subscriber) {
	p.Lock()

	//Check if the subscriber is in the subscriber list
	if _, ok := p.subs[sub]; ok {
		//If it is we close the subscriber connection
		sub.conn.Close()
		//And delete from subscriber list
		delete(p.subs, sub)
	}

	p.Unlock()
}

func (p *Publisher) broadcast(msgData []byte) error {
	for subs := range p.subs {
		fmt.Printf("Braoadcasting to %v with data %v\n", subs.conn.RemoteAddr(), msgData)
		subs.egress <- msgData
	}
	return nil
}

func (r *RoomsList) is(id int) (*Room, bool) {
	//Need to dereference the roomList to iterate over
	for room, _ := range *r {
		if room.Id == id {
			return room, true
		}
	}
	return nil, false
}
