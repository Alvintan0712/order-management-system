package main

import (
	"log"
	"net/http"

	"example.com/oms/common"
	_ "github.com/joho/godotenv/autoload"
)

var (
	addr = common.EnvString("HTTP_ADDR", ":8080")
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
