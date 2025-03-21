package xredis

import (
	"context"
	"product-store/pkg/api"
)

func (c *Client) GetProductCategory(ctx context.Context, categoryName string) (ProductCategory, error) {
	return getProductCategory(ctx, c, categoryName)
}

func (c *Client) PutProductCategory(ctx context.Context, pc api.ProductCategory) error {
	return putProductCategory(ctx, c, pc)
}

func (c *Client) GetProduct(ctx context.Context, productName string) (api.Product, error) {
	return getProduct(ctx, c, productName)
}

func (c *Client) PutProduct(ctx context.Context, p api.Product) error {
	return putProduct(ctx, c, p)
}
