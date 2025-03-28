package xredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"product-store/pkg/db"
	"product-store/pkg/types"
)

type Client struct {
	redis.UniversalClient
	Logger *zerolog.Logger
}

var _ db.DB = (*Client)(nil)

func NewDB(client redis.UniversalClient, logger *zerolog.Logger) *Client {
	c := Client{
		UniversalClient: client,
		Logger:          logger,
	}
	return &c
}

type ClientOpts struct {
	Host     string
	Port     string
	Password string
	Logger   *zerolog.Logger
}

// NewClient creates a new xredis Client with the provided connection parameters
func NewClient(opts ClientOpts) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", opts.Host, opts.Port),
		Password:    opts.Password,
		DB:          0,
		Protocol:    3,
		MaxRetries:  10,
		DialTimeout: 1 * time.Second,
	})

	client.AddHook(&loggingHook{Logger: opts.Logger})

	return &Client{
		UniversalClient: client,
		Logger:          opts.Logger,
	}
}

func (c *Client) GetProductCategory(ctx context.Context, categoryName string) (types.ProductCategory, error) {
	return getProductCategory(ctx, c, categoryName)
}

func (c *Client) PutProductCategory(ctx context.Context, pc types.ProductCategory) (types.ProductCategory, error) {
	return putProductCategory(ctx, c, pc)
}

func (c *Client) GetProduct(ctx context.Context, productName string) (types.Product, error) {
	return getProduct(ctx, c, productName)
}

func (c *Client) PutProduct(ctx context.Context, p types.Product) (types.Product, error) {
	return putProduct(ctx, c, p)
}

func (c *Client) CheckHealth(ctx context.Context) error {
	err := c.Ping(ctx).Err()
	return err
}
