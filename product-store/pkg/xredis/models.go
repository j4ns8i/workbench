package xredis

import "product-store/pkg/api"

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
	Price    float64 `redis:"price"`    // The price of the product
}

// FromAPIProduct converts a Product to its Redis representation.
func FromAPIProduct(product api.Product) Product {
	return Product{
		ID:       product.ID.String(),
		Name:     product.Name,
		Category: product.Category,
		Price:    product.Price,
	}
}

// ToAPIProduct converts a ProductRedis back to a Product.
func ToAPIProduct(obj Product) (api.Product, error) {
	id, err := api.NewULIDFromString(obj.ID)
	if err != nil {
		return api.Product{}, err
	}

	p := api.Product{
		ID: id,
		ProductData: api.ProductData{
			Name:     obj.Name,
			Category: obj.Category,
			Price:    obj.Price,
		},
	}
	return p, nil
}

// FromAPIProductCategory converts a ProductCategory to its Redis
// representation.
func FromAPIProductCategory(category api.ProductCategory) ProductCategory {
	return ProductCategory{
		ID:   category.ID.String(),
		Name: category.Name,
	}
}

// ToAPIProductCategory converts a ProductCategoryRedis back to
// a ProductCategory.
func ToAPIProductCategory(obj ProductCategory) (api.ProductCategory, error) {
	id, err := api.NewULIDFromString(obj.ID)
	if err != nil {
		return api.ProductCategory{}, err
	}

	p := api.ProductCategory{
		ID: id,
		ProductCategoryData: api.ProductCategoryData{
			Name: obj.Name,
		},
	}
	return p, nil
}
