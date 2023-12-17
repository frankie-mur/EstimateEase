package main

import (
	"estimate-ease/internal/server"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
)

type Application struct {
	port       int
	logger     *slog.Logger
	upgrader   *websocket.Upgrader
	rooms      server.RoomsList
	sync.Mutex // guards
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
	a.Lock()
	defer a.Unlock()
	a.rooms[room] = true
}

func (a *Application) removeRoom(room *server.Room) {
	a.Lock()
	defer a.Unlock()

	//Check if room exists
	if _, ok := a.rooms[room]; ok {
		delete(a.rooms, room)
	}

}
