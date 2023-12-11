package server

import (
	"sync"
)

type Publisher struct {
	//List of all subscriptions to a producer
	subs SubscriberList
	//Used to lock before editing subs
	//As maps are not concurent safe
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
		//And deletet from subscriber list
		delete(p.subs, sub)
	}

	p.Unlock()
}
