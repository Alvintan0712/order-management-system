package main

import "context"

type MenuService interface {
	CreateMenuItem(context.Context)
}

type MenuRepository interface {
	Create(context.Context) error
}
