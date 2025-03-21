package xredis

import (
	"context"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/redis/go-redis/v9"

	"product-store/pkg/api"
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

func putProduct(ctx context.Context, c redis.HashCmdable, p api.Product) error {
	key := buildProductKey(p.Name)
	obj := FromAPIProduct(p)
	return c.HMSet(ctx, key, obj).Err()
}

func getProduct(ctx context.Context, c redis.HashCmdable, productName string) (api.Product, error) {
	key := buildProductKey(productName)

	var (
		p   api.Product
		obj Product
	)
	found, err := HGetAllScan(ctx, c, key, &obj)
	if err != nil {
		return p, err
	}

	if !found {
		return p, ErrorNotFound
	}
	return ToAPIProduct(obj)
}

func putProductCategory(ctx context.Context, c redis.HashCmdable, pc api.ProductCategory) error {
	key := buildProductCategoryKey(pc.Name)
	obj := FromAPIProductCategory(pc)
	return c.HMSet(ctx, key, obj).Err()
}

func getProductCategory(ctx context.Context, c redis.HashCmdable, categoryName string) (ProductCategory, error) {
	key := buildProductCategoryKey(categoryName)

	var obj ProductCategory
	found, err := HGetAllScan(ctx, c, key, &obj)
	if err != nil {
		return obj, err
	}

	if !found {
		return obj, ErrorNotFound
	}

	return obj, nil
}
