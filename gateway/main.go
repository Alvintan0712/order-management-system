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
	httpAddr         = common.EnvString("HTTP_ADDR", ":8080")
	orderServiceAddr = common.EnvString("ORDER_SERVICE_ADDR", "localhost:8081")
)

func main() {
	var opts []grpc.DialOption

	// for tcp only connection
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(orderServiceAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Println("Dialing order service at", orderServiceAddr)

	orderClient := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(orderClient)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server:", err)
	}
}
