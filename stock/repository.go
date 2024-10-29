package main

import (
	"context"

	pb "example.com/oms/common/api"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository() (*repository, error) {
	// TODO: add database config in here
	db, err := sqlx.Connect("sqlite3", "./database/stock.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

func (r *repository) Create(ctx context.Context, entity *pb.Stock) error {
	return nil
}

func (r *repository) GetById(ctx context.Context, id string) (*pb.Stock, error) {
	return nil, nil
}

func (r *repository) Update(ctx context.Context, entity *pb.Stock) error {
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *repository) List(ctx context.Context) ([]*pb.Stock, error) {
	return nil, nil
}
