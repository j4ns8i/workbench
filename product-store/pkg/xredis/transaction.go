package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Transaction represents a Redis transaction with middleware capabilities
type Transaction struct {
	client *Client
	opts   *TransactionOptions
}

// NewTransaction creates a new Transaction with the given Redis client
func NewTransaction(client *Client) *Transaction {
	return &Transaction{
		client: client,
		opts:   &TransactionOptions{},
	}
}

// Prepare sets up the transaction with the provided options
func (t *Transaction) Prepare(options ...TransactionOption) *Transaction {
	for _, option := range options {
		option(t.opts)
	}
	return t
}

// TransactionOption represents a configuration option for a Redis transaction
type TransactionOption func(*TransactionOptions)

// TransactionOptions stores configuration for a Redis transaction
type TransactionOptions struct {
	// Keys to be watched
	watchKeys []string
	// Pre-transaction setup functions
	beforeExec []func(context.Context, *redis.Tx) error
}

// WithExists creates a transaction option that ensures the specified key
// exists before execution. This check is done under a WATCH to ensure
// consistency.
func withExists(key string, missingError error) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.watchKeys = append(opts.watchKeys, key)

		opts.beforeExec = append(opts.beforeExec, func(ctx context.Context, tx *redis.Tx) error {
			exists, err := tx.Exists(ctx, key).Result()
			if err != nil {
				return err
			}
			if exists != 1 {
				return missingError
			}
			return nil
		})
	}
}

// WithProductCategoryExists creates a transaction option that watches
// a product category and ensures it exists before executing the transaction
func WithProductCategoryExists(category string) TransactionOption {
	return func(opts *TransactionOptions) {
		key := buildProductCategoryKey(category)
		withExists(key, ErrorProductCategoryNotFound)(opts)
	}
}

// Exec executes a Redis transaction with the configured options
func (t *Transaction) Exec(ctx context.Context, fn func(context.Context, *Tx) error) error {
	return t.client.Watch(ctx, func(tx *redis.Tx) error {
		// Run all pre-transaction functions
		for _, beforeFn := range t.opts.beforeExec {
			if err := beforeFn(ctx, tx); err != nil {
				return err
			}
		}

		// Execute the user-defined function
		if err := fn(ctx, (*Tx)(tx)); err != nil {
			return err
		}

		return nil
	}, t.opts.watchKeys...)
}

type Tx redis.Tx

func (t *Tx) GetProductCategory(ctx context.Context, categoryName string) (ProductCategory, error) {
	return getProductCategory(ctx, t, categoryName)
}

func (t *Tx) GetProduct(ctx context.Context, productName string) (Product, error) {
	return getProduct(ctx, t, productName)
}

func (t *Tx) PutProductCategory(ctx context.Context, obj ProductCategory) error {
	return putProductCategory(ctx, t, obj)
}

func (t *Tx) PutProduct(ctx context.Context, obj Product) error {
	return putProduct(ctx, t, obj)
}
