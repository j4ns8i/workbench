package main

import (
	// "encoding"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewULID() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0))
}

// A fictional product
type Product struct {
	ID ulid.ULID // The unique identifier for the product
	ProductData
}

// The product data that is sent in a POST request
type ProductData struct {
	Name     string  // The name of the product
	Category string  // The name of the product category
	Price    float64 // The price of the product
}

// A fictional product category
type ProductCategory struct {
	ID                  ulid.ULID // The unique identifier for the product category
	ProductCategoryData           // The product category data
}

// The product category data that is sent in a POST request
type ProductCategoryData struct {
	Name string // The name of the product category
}

type ProductCategoryRedis struct {
	ID   string `redis:"id"`   // The unique identifier for the product category
	Name string `redis:"name"` // The name of the product category
}

func RedisFromProductCategory(category ProductCategory) ProductCategoryRedis {
	return ProductCategoryRedis{
		ID:   category.ID.String(),
		Name: category.Name,
	}
}

func NewULIDFromString(s string) (ulid.ULID, error) {
	u, err := ulid.ParseStrict(s)
	if err != nil {
		return ulid.ULID{}, err
	}
	return u, nil
}
