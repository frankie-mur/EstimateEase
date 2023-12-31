package server

import (
	"sync"
)

// Publisher manages a set of Subscribers and broadcasts messages to them.
// It uses a RWMutex to allow concurrent reads from Subscribers
// while protecting data during writes.
type Publisher struct {
	//List of all subscriptions to a producer
	Subs      SubscriberList
	callbacks []SubRemovedCallback
	//Used to lock before editing subs
	//As maps are not concurrent safe
	sync.RWMutex
}

// NewPublisher creates a new Publisher instance.
// It initializes the Subs field to an empty SubscriberList.
func NewPublisher() *Publisher {
	return &Publisher{
		Subs: make(SubscriberList),
	}
}

// AddSubscriber adds the given Subscriber to the Publisher's subscriber list.
// It acquires a write lock on the Publisher before modifying the list.
func (p *Publisher) AddSubscriber(sub *Subscriber) {
	p.Lock()
	defer p.Unlock()
	p.Subs[sub] = true
}

// RemoveSubscriber removes the given Subscriber from the Publisher's
// subscriber list. It acquires a write lock, checks if the subscriber
// is in the list, closes the subscriber's connection if found, and
// deletes the subscriber from the list.
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

	// Notify all registered callbacks
	for _, callback := range p.callbacks {
		callback.OnSubRemoved(sub)
	}
}

// Broadcast publishes the given data string to all subscribers.
// It acquires a write lock on the Publisher to synchronize access to the
// subscriber list during the broadcast. For each subscriber, it sends
// the data string over the subscriber's egress channel.
func (p *Publisher) Broadcast(data string) {
	p.Lock()
	defer p.Unlock()

	for subs := range p.Subs {
		subs.egress <- []byte(data)
	}
}

func (p *Publisher) AddCallback(callback SubRemovedCallback) {
	p.callbacks = append(p.callbacks, callback)
}
