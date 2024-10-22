package main

import (
	"log"
	"net/http"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr             = common.EnvString("HTTP_ADDR", ":8080")
	orderServiceAddr = common.EnvString("ORDER_SERVICE_ADDR", "localhost:8081")
)

func main() {
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Println("Dialing order service at", orderServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Failed to start http server:", err)
	}
}
