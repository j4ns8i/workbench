package db

import (
	"context"

	"product-store/pkg/types"
)

type DB interface {
	GetProduct(context.Context, string) (types.Product, error)
	PutProduct(context.Context, types.Product) (types.Product, error)
	GetProductCategory(context.Context, string) (types.ProductCategory, error)
	PutProductCategory(context.Context, types.ProductCategory) (types.ProductCategory, error)
	CheckHealth(context.Context) error
}
