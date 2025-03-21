package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

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
