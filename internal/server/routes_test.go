package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRoom(t *testing.T) {
	s := &Server{rooms: make(RoomsList)}
	server := httptest.NewServer(http.HandlerFunc(s.createRoom))
	defer server.Close()

	t.Run("valid request", func(t *testing.T) {
		// Arrange
		reqBody := `{"roomName":"Test Room"}`
		resp, err := http.Post(server.URL, "application/json", bytes.NewReader([]byte(reqBody)))
		if err != nil {
			t.Fatalf("error making request to server. Err: %v", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatalf("error reading response body. Err: %v", err)
		}

		// Assert
		var res Room
		err = json.Unmarshal(body, &res)
		if err != nil {
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.StatusCode)
		}

		if res.RoomName != "Test Room" {
			t.Errorf("Expected room name 'Test Room', got '%v'", res.RoomName)
		}
	})

	t.Run("missing room name", func(t *testing.T) {
		// Arrange
		s := &Server{rooms: make(RoomsList)}
		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/room", bytes.NewBufferString(reqBody))

		// Act
		w := httptest.NewRecorder()
		handler := http.HandlerFunc(s.createRoom)
		handler.ServeHTTP(w, req)

		// Assert
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest, got %v", w.Code)
		}
	})

}

//func TestJoinRoom(t *testing.T) {
//	s := &Server{rooms: make(RoomsList)}
//
//	server := httptest.NewServer(s.RegisterRoutes())
//	defer server.Close()
//
//	t.Run("valid room id", func(t *testing.T) {
//		// Arrange
//		room := NewRoom("test")
//		s.addRoom(room)
//		fmt.Printf("Room id: %d\n", room.Id)
//		url := fmt.Sprintf("%v/room/%d", server.URL, room.Id)
//		fmt.Printf("url is %v\n", url)
//		resp, err := http.Get(url)
//		if err != nil {
//			t.Fatalf("error making request to server. Err: %v", err)
//		}
//
//		defer resp.Body.Close()
//
//		if resp.StatusCode != http.StatusOK {
//			t.Errorf("Expected status OK, got %v", resp.StatusCode)
//		}
//
//	})
//
//	t.Run("invalid room id", func(t *testing.T) {
//		// Arrange
//		s := &Server{rooms: make(RoomsList)}
//		s.RegisterRoutes()
//
//		req := httptest.NewRequest(http.MethodGet, "/room/2", nil)
//
//		// Act
//		w := httptest.NewRecorder()
//		handler := http.HandlerFunc(s.joinRoom)
//		handler.ServeHTTP(w, req)
//
//		// Assert
//		if w.Code != http.StatusNotFound {
//			t.Errorf("Expected status NotFound, got %v", w.Code)
//		}
//	})

//}
