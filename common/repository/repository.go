package repository

import "context"

type Repository[T any] interface {
	Create(ctx context.Context, entity T) error
	GetById(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]T, error)
}
