package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Post("/room", s.createRoom)
	r.Get("/room/{id}", s.joinRoom)

	return r
}

// Create a new room using input name from request body
func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RoomName string `json:"roomName"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)

	//TODO: validate room name is present in request
	if input.RoomName == "" {
		http.Error(w, "Invalid room name", http.StatusBadRequest)
	}
	//create a new room
	room := NewRoom(input.RoomName)

	//add to room to the server list of rooms
	s.addRoom(room)

	fmt.Printf("Successfully created room %v", room)
	data, _ := json.Marshal(room)
	_, _ = w.Write(data)
}

// Joins a client to a room, this will be using websocket connection
func (s *Server) joinRoom(w http.ResponseWriter, r *http.Request) {
	//Get the query parameter
	id := chi.URLParam(r, "id")

	if id == "" {
		//Throw an error
		fmt.Printf("Param not specified")
		return
	}

	//Convert id to int
	i, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Id must be a number")
		return
	}

	//Check that the room exists
	room, ok := s.rooms.is(i)
	if !ok {
		http.Error(w, "Room does not exist", http.StatusNotFound)
		return
	}

	//Upgrade the websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		//TODO: use a logger here instead
		fmt.Printf("An error has occurred: %v", err)
		return
	}

	//Create a new subscriber to join
	sub := NewSubscriber(conn, room.Pub)

	//add the subscriber to the room
	room.Pub.addSubscriber(sub)

	//Start subscriber background processes to
	//read and write messages
	go sub.readMessage()
	go sub.writeMessages()

}
