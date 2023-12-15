package server

import (
	"fmt"
	"sync"
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

func (p *Publisher) Broadcast(data string) {
	for subs, _ := range p.Subs {
		fmt.Printf("Braoadcasting to %v with data %v\n", subs.conn.RemoteAddr(), data)
		subs.egress <- []byte(data)
	}
}
