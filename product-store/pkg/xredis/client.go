package xredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"product-store/pkg/db"
	"product-store/pkg/types"
)

const (
	productEventsChannelName = "productEvents"
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

func (c *Client) GetProductEvents(ctx context.Context) (<-chan types.Product, error) {
	// Create a channel for the consumer
	ch := make(chan types.Product)

	// Launch a goroutine
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		pubsub := c.Subscribe(ctx, productEventsChannelName)
		msgsCh := pubsub.Channel()

		for msg := range msgsCh {
			var p types.Product
			err := json.Unmarshal([]byte(msg.Payload), &p)
			if err != nil {
				return err
			}
			ch <- p
		}

		return pubsub.Close()
	})

	// Return the channel
	return ch, nil
}
