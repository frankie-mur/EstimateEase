package main

import (
	"estimate-ease/internal/server"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
)

type Application struct {
	port     int
	logger   *slog.Logger
	upgrader *websocket.Upgrader
	roomList *server.RoomList
}

func newApplication(newApp *Application) *http.Server {
	// Declare Server config
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", newApp.port),
		Handler:      newApp.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv
}

func (a *Application) addRoom(room *server.Room) {
	a.roomList.Lock()
	defer a.roomList.Unlock()
	a.roomList.Rooms[room] = true
}

func (a *Application) removeRoom(room *server.Room) {
	a.roomList.Lock()
	defer a.roomList.Unlock()

	//Check if room exists
	if _, ok := a.roomList.Rooms[room]; ok {
		delete(a.roomList.Rooms, room)
	}
}
