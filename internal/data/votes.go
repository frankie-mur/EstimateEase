package data

import "sync"

// Votes Represents the object to store user votes
type Votes struct {
	VoteMap map[string]string
	sync.Mutex
}

func NewVoteMap() *Votes {
	return &Votes{
		VoteMap: make(map[string]string),
	}
}

func (v *Votes) Update(name string, vote string) {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()

	v.VoteMap[name] = vote
}
