package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"product-store/pkg/types"
)

// ProductCategory is the Redis representation of a ProductCategory.
type ProductCategory struct {
	ID   string `redis:"id"`   // The unique identifier for the product category
	Name string `redis:"name"` // The name of the product category
}

// Product is the Redis representation of a Product.
type Product struct {
	ID       string  `redis:"id"`       // The unique identifier for the product
	Name     string  `redis:"name"`     // The name of the product
	Category string  `redis:"category"` // The name of the product category
	Price    float32 `redis:"price"`    // The price of the product
}

// FromAPIProduct converts a Product to its Redis representation.
func FromAPIProduct(product types.Product) Product {
	return Product{
		ID:       *product.ID,
		Name:     product.Name,
		Category: product.Category,
		Price:    product.Price,
	}
}

// ToAPIProduct converts a ProductRedis back to a Product.
func ToAPIProduct(obj Product) types.Product {
	p := types.Product{
		ID:       &obj.ID,
		Name:     obj.Name,
		Category: obj.Category,
		Price:    obj.Price,
	}
	return p
}

// FromAPIProductCategory converts a ProductCategory to its Redis
// representation.
func FromAPIProductCategory(category types.ProductCategory) ProductCategory {
	return ProductCategory{
		ID:   *category.ID,
		Name: category.Name,
	}
}

// ToAPIProductCategory converts a ProductCategoryRedis back to
// a ProductCategory.
func ToAPIProductCategory(obj ProductCategory) types.ProductCategory {
	p := types.ProductCategory{
		ID:   &obj.ID,
		Name: obj.Name,
	}
	return p
}

// HGetAllScan is a helper function to run HGETALL on a key and store the
// result in dest. dest must be an assignable address. The returned bool
// indicates if the key was found in the database.
func HGetAllScan(ctx context.Context, client redis.HashCmdable, key string, dest any) (bool, error) {
	cmd := client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	if len(cmd.Val()) == 0 {
		return false, nil
	}
	err := cmd.Scan(dest)
	if err != nil {
		return false, err
	}
	return true, nil
}
