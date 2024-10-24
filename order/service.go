package main

import (
	"context"
	"log"

	"example.com/oms/common"
	pb "example.com/oms/common/api"
)

type service struct {
	repository OrderRepository
}

func NewService(repository OrderRepository) *service {
	return &service{repository}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, r *pb.CreateOrderRequest) error {
	if len(r.Items) == 0 {
		return common.ErrNoItems
	}

	items := mergeItemsQuantities(r.Items)
	log.Println(items)

	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.Id == item.Id {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
