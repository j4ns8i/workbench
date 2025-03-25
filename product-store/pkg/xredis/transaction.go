package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Transaction represents a Redis transaction with middleware capabilities
type Transaction struct {
	client redis.UniversalClient
	opts   *TransactionOptions
}

// NewTransaction creates a new Transaction with the given Redis client
func NewTransaction(client redis.UniversalClient) *Transaction {
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
