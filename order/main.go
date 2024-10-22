package main

import "context"

func main() {
	repository := NewRepository()
	service := NewService(repository)

	service.CreateOrder(context.Background())
}
