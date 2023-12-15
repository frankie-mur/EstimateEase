package main

import (
	"estimate-ease/internal/server"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	newServer := &Server{
		port:   port,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		rooms: make(server.RoomsList),
	}

	srv := NewServer(newServer)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
