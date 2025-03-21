package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"product-store/pkg/store"
	"product-store/pkg/types"
)

type Tx redis.Tx

var _ store.Store = (*Tx)(nil)

func (t *Tx) GetProductCategory(ctx context.Context, name string) (types.ProductCategory, error) {
	return getProductCategory(ctx, t, name)
}

func (t *Tx) GetProduct(ctx context.Context, name string) (types.Product, error) {
	return getProduct(ctx, t, name)
}

func (t *Tx) PutProductCategory(ctx context.Context, obj types.ProductCategory) error {
	return putProductCategory(ctx, t, obj)
}

func (t *Tx) PutProduct(ctx context.Context, p types.Product) error {
	return putProduct(ctx, t, p)
}
