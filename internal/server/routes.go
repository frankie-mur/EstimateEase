package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/ws", s.HelloWebSocketHandler)

	return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) HelloWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		//Internal server error
		fmt.Printf("An error has occurred: %v", err)
		return
	}

	//Create a new producer
	pub := NewPublisher()

	sub := NewSubscriber(conn, pub)

	pub.addSubscriber(NewSubscriber(conn, pub))
	//Start subscriper background proccesses to
	//read and write messages
	go sub.readMessage()

}
