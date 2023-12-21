package main

import (
	"estimate-ease/internal/server"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	app := &Application{
		port:   port,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		roomList: server.NewRoomList(),
	}

	srv := newApplication(app)

	//goroutine that check for idle rooms and remove them
	//for a given time interval
	go app.checkIdleRooms(5 * time.Minute)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
