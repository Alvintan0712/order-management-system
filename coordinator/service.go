package main

type CoordinatorService interface {
}

type service struct {
}

func NewService() *service {
	return &service{}
}
