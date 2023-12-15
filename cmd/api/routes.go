package main

import (
	"errors"
	"estimate-ease/internal/server"
	"estimate-ease/ui/components"
	"estimate-ease/ui/data"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	r.Get("/", s.homePage)
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

	err := s.writeJSON(w, http.StatusCreated, room.Id, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

}

// Joins a client to a room, this will be using websocket connection
func (s *Server) joinRoom(w http.ResponseWriter, r *http.Request) {
	//Get the query parameter
	id := r.FormValue("roomID")
	displayName := r.FormValue("displayName")

	if id == "" {
		s.badRequestResponse(w, r, errors.New("invalid room ID"))
		return
	}

	if displayName == "" {
		s.badRequestResponse(w, r, errors.New("invalid display name"))
		return
	}

	//Check that the room exists
	_, ok := s.rooms.Is(id)
	if !ok {
		s.roomDoesNotExistResponse(w, r)
		return
	}

	url := fmt.Sprintf("/room/%v/%v", id, displayName)
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) connectToRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")
	displayName := chi.URLParam(r, "displayName")

	if id == "" {
		s.badRequestResponse(w, r, errors.New("invalid room ID"))
		return
	}

	if displayName == "" {
		s.badRequestResponse(w, r, errors.New("invalid display name"))
		return
	}

	//Check that the room exists
	room, ok := s.rooms.Is(id)
	if !ok {
		s.roomDoesNotExistResponse(w, r)
		return
	}
	//Upgrade the websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

	//Create a new subscriber to join
	sub := server.NewSubscriber(conn, room, displayName)

	//add the subscriber to the room
	room.Pub.AddSubscriber(sub)

	//Update user count for all subscribers
	htmlResponse := fmt.Sprintf("<div id=\"room-count\">%d</div>", len(sub.Publisher.Subs))
	go sub.Publisher.Broadcast(htmlResponse)
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	c := components.HomePage()
	err := c.Render(r.Context(), w)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}
}

func (s *Server) roomPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")
	displayName := chi.URLParam(r, "displayName")

	if id == "" {
		s.badRequestResponse(w, r, errors.New("invalid room ID"))
		return
	}

	if displayName == "" {
		s.badRequestResponse(w, r, errors.New("invalid display name"))
		return
	}

	//Check that the room exists
	room, ok := s.rooms.Is(id)
	if !ok {
		s.roomDoesNotExistResponse(w, r)
		return
	}

	pageData := data.RoomPageData{
		Room:        room,
		DisplayName: displayName,
	}

	c := components.RoomPage(pageData)
	err := c.Render(r.Context(), w)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
