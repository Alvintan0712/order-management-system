package main

import (
	"context"
	"fmt"

	"example.com/oms/common/api/avro"
	pb "example.com/oms/common/api/protobuf"
	"example.com/oms/common/broker/producer"
	"example.com/oms/common/broker/topic"
	rp "example.com/oms/common/repository"
	"example.com/oms/common/schemaregistry/serializer"
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
	repository rp.Repository[*pb.MenuItem]
	producer   producer.Producer
	serializer serializer.Serializer
}

func NewService(repository rp.Repository[*pb.MenuItem], producer producer.Producer, serializer serializer.Serializer) *service {
	return &service{
		repository: repository,
		producer:   producer,
		serializer: serializer,
	}
}

func (s *service) CreateMenuItem(ctx context.Context, r *pb.CreateMenuItemRequest) (*pb.MenuItem, error) {
	menu := &pb.MenuItem{
		Id:        uuid.New().String(),
		Name:      r.Name,
		UnitPrice: r.UnitPrice,
		Currency:  r.Currency,
	}
	err := s.repository.Create(ctx, menu)
	if err != nil {
		return nil, fmt.Errorf("failed to create entity: %v", err)
	}

	event := &avro.MenuCreatedEvent{
		Id:        menu.Id,
		Name:      menu.Name,
		UnitPrice: int(menu.UnitPrice),
		Currency:  menu.Currency,
	}
	payload, err := s.serializer.Serialize(topic.MenuCreated, event)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %v", err)
	}

	err = s.producer.Produce(topic.MenuCreated, []byte(menu.Id), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to produce event: %v", err)
	}

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
