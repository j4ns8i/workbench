package xredis

import (
	"context"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/redis/go-redis/v9"

	"product-store/pkg/types"
)

func buildRedisKey(kind, name string) string {
	hexDigest := fmt.Sprintf("%x", xxhash.Sum64String(name))
	return kind + ":" + hexDigest
}

func buildProductKey(productName string) string {
	return buildRedisKey("PRODUCT", productName)
}

func buildProductCategoryKey(categoryName string) string {
	return buildRedisKey("PRODUCTCATEGORY", categoryName)
}

func putProduct(ctx context.Context, c redis.HashCmdable, p types.Product) error {
	key := buildProductKey(p.Name)
	obj := FromAPIProduct(p)
	return c.HMSet(ctx, key, obj).Err()
}

func getProduct(ctx context.Context, c redis.HashCmdable, productName string) (types.Product, error) {
	key := buildProductKey(productName)

	var (
		p   types.Product
		obj Product
	)
	found, err := HGetAllScan(ctx, c, key, &obj)
	if err != nil {
		return p, err
	}

	if !found {
		return p, ErrNotFound
	}
	return ToAPIProduct(obj)
}

func putProductCategory(ctx context.Context, c redis.HashCmdable, pc types.ProductCategory) error {
	key := buildProductCategoryKey(pc.Name)
	obj := FromAPIProductCategory(pc)
	return c.HMSet(ctx, key, obj).Err()
}

func getProductCategory(ctx context.Context, c redis.HashCmdable, categoryName string) (types.ProductCategory, error) {
	key := buildProductCategoryKey(categoryName)

	var (
		pc  types.ProductCategory
		obj ProductCategory
	)
	found, err := HGetAllScan(ctx, c, key, &obj)
	if err != nil {
		return pc, err
	}

	if !found {
		return pc, ErrNotFound
	}

	return ToAPIProductCategory(obj)
}
