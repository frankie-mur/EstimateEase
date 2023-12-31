package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/frankie-mur/EstimateEase/internal/server"
	"github.com/gorilla/sessions"
)

func TestCreateRoom(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room", nil)
		req.Form = url.Values{"roomName": {"Test Room"}}

		w := httptest.NewRecorder()
		a := &Application{roomList: server.NewRoomList(), sessionStore: sessions.NewCookieStore([]byte("test"))}

		a.createRoom(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("Expected status Created but got %v", w.Code)
		}
	})

	t.Run("empty room name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/room", nil)

		w := httptest.NewRecorder()
		a := &Application{roomList: server.NewRoomList()}
		a.createRoom(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest but got %v", w.Code)
		}
	})

}

func TestJoinRoom(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/room/join", nil)
	a := &Application{roomList: server.NewRoomList()}

	t.Run("valid room ID and display name", func(t *testing.T) {
		w := httptest.NewRecorder()

		newRoom := server.NewRoom("Test Room")
		roomId := newRoom.Id
		a.addRoom(newRoom)
		req.Form = url.Values{"roomID": {roomId},
			"displayName": {"John"}}

		a.joinRoom(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("Expected redirect status but got %v", w.Code)
		}
	})

	t.Run("missing room ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req.Form = url.Values{"displayName": {"John"}}

		a.joinRoom(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected bad request status but got %v", w.Code)
		}
	})

	t.Run("invalid room ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req.Form = url.Values{"roomID": {"invalid"},
			"displayName": {"John"}}

		a.joinRoom(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected bad request status but got %v", w.Code)
		}
	})

	t.Run("missing display name", func(t *testing.T) {
		w := httptest.NewRecorder()
		req.Form = url.Values{"roomID": {"123e4567-e89b-12d3-a456-426614174000"}}

		a.joinRoom(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected bad request status but got %v", w.Code)
		}
	})

	t.Run("room not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req.Form = url.Values{"roomID": {"123e4567-e89b-12d3-a456-426614174000"},
			"displayName": {"John"}}

		a.joinRoom(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected not found status but got %v", w.Code)
		}
	})

}
