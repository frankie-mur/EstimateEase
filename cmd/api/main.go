package main

import (
	"estimate-ease/internal/server"
)

func main() {
	srv := server.NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
