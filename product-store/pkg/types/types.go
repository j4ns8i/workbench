package types

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

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

func NewULID() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
}

func NewULIDFromString(s string) (ulid.ULID, error) {
	u, err := ulid.ParseStrict(s)
	if err != nil {
		return ulid.ULID{}, err
	}
	return u, nil
}
