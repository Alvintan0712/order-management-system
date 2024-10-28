package main

import (
	"context"
	"database/sql"

	pb "example.com/oms/common/api"
	_ "github.com/mattn/go-sqlite3"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity T) error
	// GetById(ctx context.Context, id string) (T, error)
	// Update(ctx context.Context, entity T) error
	// Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]T, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository() (*repository, error) {
	// TODO: add database config in here
	db, err := sql.Open("sqlite3", "./menu.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

func (r *repository) Create(ctx context.Context, entity *pb.MenuItem) error {
	query := "INSERT INTO menu (id, name, unitPrice, currency) VALUES (?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query, entity.Id, entity.Name, entity.UnitPrice, entity.Currency)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) List(ctx context.Context) ([]*pb.MenuItem, error) {
	query := "SELECT id, name, unitPrice, currency FROM menu;"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	items := []*pb.MenuItem{}
	for rows.Next() {
		var menuItem pb.MenuItem
		if err := rows.Scan(&menuItem.Id, &menuItem.Name, &menuItem.UnitPrice, &menuItem.Currency); err != nil {
			return nil, err
		}
		items = append(items, &menuItem)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
