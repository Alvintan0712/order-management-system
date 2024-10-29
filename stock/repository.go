package main

import (
	"context"
	"database/sql"

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
	query := `
		INSERT INTO stock (item_id, quantity)
		VALUES(:item_id, :quantity)
	`

	_, err := r.db.NamedExecContext(ctx, query, entity.ToDB())
	if err != nil {
		return nil
	}

	return nil
}

func (r *repository) GetById(ctx context.Context, id string) (*pb.Stock, error) {
	query := `
		SELECT item_id, quantity
		FROM stock
		WHERE item_id = $1
		LIMIT 1
	`

	var stockDB pb.StockDB
	err := r.db.GetContext(ctx, &stockDB, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return stockDB.ToProto(), nil
}

func (r *repository) Update(ctx context.Context, entity *pb.Stock) error {
	query := `
		UPDATE stock
		SET quantity = :quantity
		WHERE item_id = :item_id
	`

	_, err := r.db.NamedExecContext(ctx, query, entity.ToDB())
	if err != nil {
		return nil
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM stock
		WHERE item_id = ?;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) List(ctx context.Context) ([]*pb.Stock, error) {
	query := `
		SELECT item_id, quantity
		FROM stock;
	`

	var dbStocks []*pb.StockDB
	err := r.db.SelectContext(ctx, &dbStocks, query)
	if err != nil {
		return nil, err
	}

	stocks := make([]*pb.Stock, len(dbStocks))
	for i, stock := range dbStocks {
		stocks[i] = stock.ToProto()
	}

	return stocks, nil
}
