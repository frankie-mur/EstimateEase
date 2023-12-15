package main

import "log"

func main() {
	srv := NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
