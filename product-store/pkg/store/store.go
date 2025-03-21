package store

import (
	"context"

	"product-store/pkg/api"
)

type Store interface {
	GetProduct(context.Context, string) (api.Product, error)
	PutProduct(context.Context, api.Product) error
	GetProductCategory(context.Context, string) (api.ProductCategory, error)
	PutProductCategory(context.Context, api.ProductCategory) error
}
