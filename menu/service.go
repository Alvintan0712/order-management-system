package main

import (
	"context"

	pb "example.com/oms/common/api"
	"github.com/google/uuid"
)

type MenuService interface {
	CreateMenuItem(context.Context, *pb.CreateMenuItemRequest) (*pb.MenuItem, error)
	GetMenuItem(context.Context, string) (*pb.MenuItem, error)
	UpdateMenuItem(context.Context, *pb.UpdateMenuItemRequest) error
	DeleteMenuItem(context.Context, string) error
	ListMenuItems(context.Context) ([]*pb.MenuItem, error)
}

type service struct {
	repository Repository[*pb.MenuItem]
}

func NewService(repository Repository[*pb.MenuItem]) *service {
	return &service{repository: repository}
}

func (s *service) CreateMenuItem(ctx context.Context, r *pb.CreateMenuItemRequest) (*pb.MenuItem, error) {
	menu := &pb.MenuItem{
		Id:        uuid.New().String(),
		Name:      r.Name,
		UnitPrice: r.UnitPrice,
		Currency:  r.Currency,
	}
	s.repository.Create(ctx, menu)
	return menu, nil
}

func (s *service) GetMenuItem(ctx context.Context, id string) (*pb.MenuItem, error) {
	return s.repository.GetById(ctx, id)
}

func (s *service) UpdateMenuItem(ctx context.Context, r *pb.UpdateMenuItemRequest) error {
	menu := &pb.MenuItem{
		Id:        r.Id,
		Name:      r.Name,
		UnitPrice: r.UnitPrice,
		Currency:  r.Currency,
	}
	return s.repository.Update(ctx, menu)
}

func (s *service) DeleteMenuItem(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) ListMenuItems(ctx context.Context) ([]*pb.MenuItem, error) {
	return s.repository.List(ctx)
}
