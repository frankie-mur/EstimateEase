package main

import (
	"estimate-ease/internal/server"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port       int
	upgrader   *websocket.Upgrader
	rooms      server.RoomsList
	sync.Mutex // guards
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		rooms: make(server.RoomsList),
	}

	// Declare Server config
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
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
