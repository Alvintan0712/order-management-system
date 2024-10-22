package main

import (
	"log"
	"net/http"
)

const (
	addr = ":8080"
)

func main() {
	mux := http.NewServeMux()
	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Failed to start http server:", err)
	}
}
