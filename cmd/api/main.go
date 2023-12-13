package main

func main() {
	srv := NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
