package server

// Room represents a room in the server (or session)
// Each room will have one publisher
// a publisher can have one or more subscribers
type Room struct {
	Id       int        `json:"id"`
	RoomName string     `json:"room_name"`
	Pub      *Publisher `json:"-"`
}

type RoomsList map[*Room]bool

func NewRoom(name string) *Room {
	//TODO: How do we want to handle the id?
	return &Room{
		Id:       0,
		RoomName: name,
		Pub:      NewPublisher(),
	}
}
