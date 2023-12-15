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

type Server struct {
	port       int
	logger     *slog.Logger
	upgrader   *websocket.Upgrader
	rooms      server.RoomsList
	sync.Mutex // guards
}

func NewServer(newServer *Server) *http.Server {
	// Declare Server config
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv
}

func (s *Server) addRoom(room *server.Room) {
	s.Lock()
	defer s.Unlock()
	s.rooms[room] = true
}

func (s *Server) removeRoom(room *server.Room) {
	s.Lock()
	defer s.Unlock()

	//Check if room exists
	if _, ok := s.rooms[room]; ok {
		delete(s.rooms, room)
	}

}
