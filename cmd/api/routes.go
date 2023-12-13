package main

import (
	"encoding/json"
	"estimate-ease/internal/server"
	"estimate-ease/ui/components"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", s.home)
	r.Post("/room", s.createRoom)
	r.Post("/room/join", s.joinRoom)
	r.Get("/ws/room/{roomID}", s.connectToRoom)
	r.Get("/room/{roomID}", s.roomPage)

	return r
}

// Create a new room using input name from request body
func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	//var input struct {
	//	RoomName string `json:"roomName"`
	//}
	roomName := r.FormValue("roomName")
	//_ = json.NewDecoder(r.Body).Decode(&input)

	if roomName == "" {
		http.Error(w, "Invalid room name", http.StatusBadRequest)
		return
	}
	//create a new room
	room := server.NewRoom(roomName)

	//add to room to the server list of rooms
	s.addRoom(room)

	fmt.Printf("Successfully created room %v\n", room)
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(room)
	_, _ = w.Write(data)
}

// Joins a client to a room, this will be using websocket connection
func (s *Server) joinRoom(w http.ResponseWriter, r *http.Request) {
	//Get the query parameter
	id := r.FormValue("roomID")

	if id == "" {
		//Throw an error
		fmt.Printf("Param not specified")
		return
	}

	uuidVal, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Error parsing uuid: %v", err)
		http.Error(w, "Invalid uuid", http.StatusBadRequest)
		return
	}

	//Check that the room exists
	_, ok := s.rooms.Is(uuidVal)
	if !ok {
		http.Error(w, "Room does not exist", http.StatusNotFound)
		return
	}

	fmt.Printf("Redirecting to room %v\n", uuidVal)
	url := fmt.Sprintf("/room/%v", uuidVal.String())
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) connectToRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")
	fmt.Printf("Connecting to room %v\n", id)

	if id == "" {
		//Throw an error
		fmt.Printf("Param not specified")
		return
	}

	uuidVal, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Error parsing uuid: %v", err)
		http.Error(w, "Invalid uuid", http.StatusBadRequest)
		return
	}

	//Check that the room exists
	room, ok := s.rooms.Is(uuidVal)
	if !ok {
		http.Error(w, "Room does not exist", http.StatusNotFound)
		return
	}
	fmt.Printf("Upgrading connection")
	//Upgrade the websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		//TODO: use a logger here instead
		fmt.Printf("An error has occurred: %v", err)
		return
	}

	//Create a new subscriber to join
	sub := server.NewSubscriber(conn, room.Pub)

	//add the subscriber to the room
	room.Pub.AddSubscriber(sub)

	//Start subscriber background processes to
	//read and write messages
	go sub.ReadMessage()
	go sub.WriteMessages()
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	c := components.HomePage("Frankie")
	err := c.Render(r.Context(), w)
	if err != nil {
		return
	}
}

func (s *Server) roomPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")

	if id == "" {
		//Throw an error
		fmt.Printf("Param not specified")
		return
	}

	uuidVal, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Error parsing uuid: %v", err)
		http.Error(w, "Invalid uuid", http.StatusBadRequest)
		return
	}

	//Check that the room exists
	room, ok := s.rooms.Is(uuidVal)
	if !ok {
		http.Error(w, "Room does not exist", http.StatusNotFound)
		return
	}
	c := components.RoomPage(*room)
	err = c.Render(r.Context(), w)
	if err != nil {
		return
	}
}
