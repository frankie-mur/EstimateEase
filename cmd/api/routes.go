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
	r.Get("/ws/room/{roomID}/{displayName}", s.connectToRoom)
	r.Get("/room/{roomID}/{displayName}", s.roomPage)

	return r
}

// Create a new room using input name from request body
func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) {
	roomName := r.FormValue("roomName")

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
	displayName := r.FormValue("displayName")

	if id == "" {
		//Throw an error
		fmt.Printf("Param id not specified")
		return
	}

	if displayName == "" {
		fmt.Printf("displayName not specified")
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
	url := fmt.Sprintf("/room/%v/%v", uuidVal.String(), displayName)
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) connectToRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")
	displayName := chi.URLParam(r, "displayName")
	fmt.Printf("Connecting to room %v\n", id)

	if id == "" {
		//Throw an error
		fmt.Printf("Param not specified")
		return
	}

	if displayName == "" {
		fmt.Printf("displayName not specified")
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
	sub := server.NewSubscriber(conn, room, displayName)

	//add the subscriber to the room
	room.Pub.AddSubscriber(sub)
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
	displayName := chi.URLParam(r, "displayName")

	if id == "" {
		//Throw an error
		fmt.Printf("Param not specified")
		return
	}

	if displayName == "" {
		fmt.Printf("displayName not specified")
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

	pageData := server.RoomPageData{
		Room:        room,
		DisplayName: displayName,
	}

	c := components.RoomPage(pageData)
	err = c.Render(r.Context(), w)
	if err != nil {
		return
	}
}
