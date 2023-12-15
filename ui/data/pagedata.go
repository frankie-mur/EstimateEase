package data

import "estimate-ease/internal/server"

// RoomPageData data needed for the room page html template
type RoomPageData struct {
	Room        *server.Room
	DisplayName string
}

type VoteMapData struct {
	VoteMap map[string]string
}
