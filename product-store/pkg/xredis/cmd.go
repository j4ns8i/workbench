package xredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"

	"product-store/pkg/db"
	"product-store/pkg/ptr"
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

func putProduct(ctx context.Context, c redis.UniversalClient, p types.Product) (types.Product, error) {
	t := NewTransaction(c)
	t.Prepare(WithProductCategoryExists(p.Category))
	var result types.Product
	err := t.Exec(ctx, func(ctx context.Context, tx *Tx) error {
		// Check if the product already exists to preserve its ID
		var exists = true
		existingObj, err := tx.GetProduct(ctx, p.Name)
		if err != nil {
			if errors.Is(err, db.ErrProductNotFound) {
				exists = false
			} else {
				return err
			}
		}
		if !exists {
			p.ID = ptr.To(ulid.Make().String())
		} else {
			p.ID = existingObj.ID
		}

		key := buildProductKey(p.Name)
		obj := FromAPIProduct(p)
		err = c.HSet(ctx, key, obj).Err()
		if err != nil {
			return err
		}
		result = p
		return nil
	})
	return result, err
}

func getProduct(ctx context.Context, c redis.HashCmdable, name string) (types.Product, error) {
	var (
		key = buildProductKey(name)
		obj Product
	)
	exists, err := HGetAllScan(ctx, c, key, &obj)
	if err != nil {
		return types.Product{}, err
	}
	if !exists {
		return types.Product{}, db.ErrProductNotFound
	}
	return ToAPIProduct(obj), nil
}

func putProductCategory(ctx context.Context, c redis.HashCmdable, pc types.ProductCategory) (types.ProductCategory, error) {
	var (
		existing ProductCategory
		key      = buildProductCategoryKey(pc.Name)
	)
	exists, err := HGetAllScan(ctx, c, key, &existing) // TODO: race condition here
	if err != nil {
		return types.ProductCategory{}, err
	}
	if !exists {
		// not found, create a new ID
		pc.ID = ptr.To(ulid.Make().String())
	} else {
		pc.ID = &existing.ID
	}

	obj := FromAPIProductCategory(pc)
	err = c.HSet(ctx, key, obj).Err()
	return pc, err
}

func getProductCategory(ctx context.Context, c redis.HashCmdable, name string) (types.ProductCategory, error) {
	var (
		key = buildProductCategoryKey(name)
		obj ProductCategory
	)
	exists, err := HGetAllScan(ctx, c, key, &obj)
	if err != nil {
		return types.ProductCategory{}, err
	}
	if !exists {
		return types.ProductCategory{}, db.ErrProductCategoryNotFound
	}
	return ToAPIProductCategory(obj), nil
}
