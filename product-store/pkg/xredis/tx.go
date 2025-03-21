package xredis

import (
	"context"
	"product-store/pkg/api"

	"github.com/redis/go-redis/v9"
)

type Tx redis.Tx

func (t *Tx) GetProductCategory(ctx context.Context, name string) (ProductCategory, error) {
	return getProductCategory(ctx, t, name)
}

func (t *Tx) GetProduct(ctx context.Context, name string) (api.Product, error) {
	return getProduct(ctx, t, name)
}

func (t *Tx) PutProductCategory(ctx context.Context, obj api.ProductCategory) error {
	return putProductCategory(ctx, t, obj)
}

func (t *Tx) PutProduct(ctx context.Context, p api.Product) error {
	return putProduct(ctx, t, p)
}
