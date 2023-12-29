package main

import (
	"bytes"
	"errors"
	"estimate-ease/internal/server"
	"estimate-ease/ui/components"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

// TOOD: make this variable dynamic
var HOST = "localhost:8080"

func (a *Application) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	r.Get("/", a.homePage)
	r.Post("/room", a.createRoom)
	r.Post("/room/join", a.joinRoom)
	r.Get("/ws/room/{roomID}/{displayName}", a.connectToRoom)
	r.Get("/room/{roomID}/{displayName}", a.roomPage)
	r.Get("/room/{roomID}", a.displayNamePage)
	r.Post("/room/join/user", a.getDisplayNameAndRoute)

	return r
}

// Create a new room using input name from request body
func (a *Application) createRoom(w http.ResponseWriter, r *http.Request) {
	roomName := r.FormValue("roomName")

	if roomName == "" {
		http.Error(w, "Invalid room name", http.StatusBadRequest)
		return
	}
	//create a new room
	room := server.NewRoom(roomName)

	//add to room to the server list of rooms
	a.addRoom(room)

	//Add created flash message to session storage
	err := a.addSessionFlash(r, w, "estimate-ease", "createdFlash")
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

	//route to room page
	http.Redirect(w, r, fmt.Sprintf("/room/%v", room.Id), http.StatusFound)

}

// Joins a client to a room, this will be using websocket connection
func (a *Application) joinRoom(w http.ResponseWriter, r *http.Request) {
	//Get the query parameter
	id := r.FormValue("roomID")
	displayName := r.FormValue("displayName")

	if id == "" {
		a.badRequestResponse(w, r, errors.New("invalid room ID"))
		return
	}

	if displayName == "" {
		a.badRequestResponse(w, r, errors.New("invalid display name"))
		return
	}

	//Check that the room exists
	_, ok := a.roomList.Is(id)
	if !ok {
		a.roomDoesNotExistResponse(w, r)
		return
	}

	url := fmt.Sprintf("/room/%v/%v", id, displayName)
	http.Redirect(w, r, url, http.StatusFound)
}

func (a *Application) connectToRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")
	displayName := chi.URLParam(r, "displayName")

	if id == "" {
		a.badRequestResponse(w, r, errors.New("invalid room ID"))
		return
	}

	if displayName == "" {
		a.badRequestResponse(w, r, errors.New("invalid display name"))
		return
	}

	//Check that the room exists
	room, ok := a.roomList.Is(id)
	if !ok {
		a.roomDoesNotExistResponse(w, r)
		return
	}
	//Upgrade the websocket connection
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	//Create a new subscriber to join
	sub := server.NewSubscriber(conn, room, displayName)

	//add the subscriber to the room
	room.Pub.AddSubscriber(sub)

	//Update user count for all subscribers
	numUsers := fmt.Sprintf("%d", (len(room.Pub.Subs)))
	var buf bytes.Buffer
	data, err := server.RenderComponentToBuffer(
		components.Stats("Total Users", numUsers),
		r.Context(),
		&buf,
	)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
	go sub.Publisher.Broadcast(data)
}

func (a *Application) homePage(w http.ResponseWriter, r *http.Request) {
	c := components.HomePage(fmt.Sprintf("%d", len(a.roomList.Rooms)))
	err := c.Render(r.Context(), w)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *Application) roomPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roomID")
	displayName := chi.URLParam(r, "displayName")

	if id == "" {
		a.badRequestResponse(w, r, errors.New("invalid room ID"))
		return
	}

	if displayName == "" {
		a.badRequestResponse(w, r, errors.New("invalid display name"))
		return
	}

	//Check that the room exists
	room, ok := a.roomList.Is(id)
	if !ok {
		a.roomDoesNotExistResponse(w, r)
		return
	}

	pageData := components.RoomPageData{
		RoomName:    room.RoomName,
		RoomID:      room.Id,
		DisplayName: displayName,
		RoomURL:     fmt.Sprintf("%v/room/%v", HOST, room.Id),
	}

	c := components.RoomPage(pageData)
	err := c.Render(r.Context(), w)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *Application) displayNamePage(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomID")
	ses, _ := a.sessionStore.Get(r, "estimate-ease")
	// Get the previous flashes, if any.
	showFlash := false
	if flashes := ses.Flashes("createdFlash"); len(flashes) > 0 {
		showFlash = true
	}
	ses.Save(r, w)
	c := components.DisplayNamePage(roomId, showFlash)
	err := c.Render(r.Context(), w)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *Application) getDisplayNameAndRoute(w http.ResponseWriter, r *http.Request) {
	displayName := r.FormValue("displayName")
	roomID := r.FormValue("roomID")

	if displayName == "" {
		a.badRequestResponse(w, r, errors.New("invalid display name"))
	}
	url := fmt.Sprintf("/room/%v/%v", roomID, displayName)
	http.Redirect(w, r, url, http.StatusFound)
}
