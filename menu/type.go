package main

import (
	"context"

	pb "example.com/oms/common/api"
)

type MenuService interface {
	CreateMenuItem(context.Context, *pb.CreateMenuItemRequest) (*pb.MenuItem, error)
}

type MenuRepository interface {
	Create(context.Context) error
}
