package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"product-store/pkg/types"
)

type Tx redis.Tx

func (t *Tx) GetProductCategory(ctx context.Context, name string) (types.ProductCategory, error) {
	return getProductCategory(ctx, t, name)
}

func (t *Tx) GetProduct(ctx context.Context, name string) (types.Product, error) {
	return getProduct(ctx, t, name)
}
