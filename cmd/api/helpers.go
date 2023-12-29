package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Helper function that will marshal data to json and write back with provided headers and status code
func (a *Application) writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// checkIdleRooms periodically checks for chat rooms that have no active
// users and removes them from the app.rooms map. It uses a ticker to
// check at a regular interval.
func (app *Application) checkIdleRooms(timeBetweenRequests time.Duration) {
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		//loop throuhgh all rooms and check if a room has zero users
		log.Printf("Checking for idle rooms...")
		for room := range app.roomList.Rooms {
			if len(room.Pub.Subs) == 0 {
				//remove room from app.rooms
				app.removeRoom(room)
				log.Printf("Room %s has no users, removing...", room.Id)
			}
		}
	}
}

func (app *Application) addSessionFlash(r *http.Request, w http.ResponseWriter, sessionName, value string) error {
	ses, err := app.sessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}
	ses.AddFlash(nil, value)
	err = ses.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}
