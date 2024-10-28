package main

import (
	"context"
	"log"
)

type service struct {
	repository MenuRepository
}

func NewService(repository MenuRepository) *service {
	return &service{repository}
}

func (s *service) CreateMenuItem(context.Context) {
	log.Println("Create menu item")
}
