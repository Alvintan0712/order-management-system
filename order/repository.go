package main

import "context"

type repository struct {
}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) Create(context.Context) error {
	return nil
}
