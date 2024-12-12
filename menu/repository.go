package main

import (
	"context"

	pb "example.com/oms/common/api/protobuf"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository() (*repository, error) {
	// TODO: add database config in here
	db, err := sqlx.Connect("sqlite3", "./database/menu.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

func (r *repository) Create(ctx context.Context, entity *pb.MenuItem) error {
	query := `
		INSERT INTO menu (id, name, unitPrice, currency)
		VALUES (:id, :name, :unitPrice, :currency);
	`

	_, err := r.db.NamedExecContext(ctx, query, entity.ToDB())
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetById(ctx context.Context, id string) (*pb.MenuItem, error) {
	query := `
		SELECT id, name, unitPrice, currency
		FROM menu
		WHERE id = $1
		LIMIT 1;
	`

	var itemDB pb.MenuItemDB
	if err := r.db.GetContext(ctx, &itemDB, query, id); err != nil {
		return nil, err
	}

	return itemDB.ToProto(), nil
}

func (r *repository) Update(ctx context.Context, entity *pb.MenuItem) error {
	query := `
		UPDATE menu
		SET name=:name, unitPrice=:unitPrice, currency=:currency
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, entity.ToDB())
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM menu
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) List(ctx context.Context) ([]*pb.MenuItem, error) {
	query := `
		SELECT id, name, unitPrice, currency
		FROM menu;
	`

	var dbItems []*pb.MenuItemDB
	err := r.db.SelectContext(ctx, &dbItems, query)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.MenuItem, len(dbItems))
	for i, item := range dbItems {
		items[i] = item.ToProto()
	}

	return items, nil
}
