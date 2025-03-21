package store

import (
	"context"

	"product-store/pkg/types"
)

type Store interface {
	GetProduct(context.Context, string) (types.Product, error)
	PutProduct(context.Context, types.Product) error
	GetProductCategory(context.Context, string) (types.ProductCategory, error)
	PutProductCategory(context.Context, types.ProductCategory) error
}
