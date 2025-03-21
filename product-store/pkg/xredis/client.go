package xredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"product-store/pkg/api"
	"product-store/pkg/store"
)

type Client struct {
	redis.UniversalClient
	Logger *zerolog.Logger
}

var _ store.Store = (*Client)(nil)

type ClientOpts struct {
	Host     string
	Port     string
	Password string
	Logger   *zerolog.Logger
}

// NewClient creates a new xredis Client with the provided connection parameters
func NewClient(opts ClientOpts) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", opts.Host, opts.Port),
		Password: opts.Password,
		DB:       0,
		Protocol: 3,
	})

	client.AddHook(&loggingHook{Logger: opts.Logger})

	return &Client{
		UniversalClient: client,
		Logger:          opts.Logger,
	}
}

func (c *Client) GetProductCategory(ctx context.Context, categoryName string) (api.ProductCategory, error) {
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
