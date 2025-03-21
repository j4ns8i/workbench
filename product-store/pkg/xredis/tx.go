package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"product-store/pkg/api"
	"product-store/pkg/store"
)

type Tx redis.Tx

var _ store.Store = (*Tx)(nil)

func (t *Tx) GetProductCategory(ctx context.Context, name string) (api.ProductCategory, error) {
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
