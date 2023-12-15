package main

import (
	"estimate-ease/internal/server"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateRoom(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room", nil)
		req.Form = url.Values{"roomName": {"Test Room"}}

		w := httptest.NewRecorder()
		s := &Server{rooms: make(server.RoomsList)}

		s.createRoom(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status Created but got %v", w.Code)
		}
	})

	t.Run("empty room name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room", nil)

		w := httptest.NewRecorder()
		s := &Server{rooms: make(server.RoomsList)}
		s.createRoom(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest but got %v", w.Code)
		}
	})

}
