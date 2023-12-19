package server

import (
	"sort"
	"sync"
)

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

// Function reurn list of sorted names alphabetically
func (v *Votes) sortNames() []string {
	v.Mutex.Lock()
	defer v.Mutex.Unlock()

	// Extract keys from map
	sortedNames := make([]string, 0, len(v.VoteMap))
	for k := range v.VoteMap {
		sortedNames = append(sortedNames, k)
	}
	// Sort keys
	sort.Strings(sortedNames)

	return sortedNames
}
