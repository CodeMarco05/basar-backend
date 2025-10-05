package main

import (
	"hackathon-basar-backend/internal/server"
	"log"
	"net/http"
)

func main() {
	srv := server.NewServer()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
