package main

import "context"

type service struct {
	repository OrderRepository
}

func NewService(repository OrderRepository) *service {
	return &service{repository}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}
