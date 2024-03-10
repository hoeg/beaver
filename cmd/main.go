package main

import (
	"log"
	"net/http"

	"github.com/hoeg/beaver/internal/server"
)

func main() {
	srv := server.New(nil)

	log.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
