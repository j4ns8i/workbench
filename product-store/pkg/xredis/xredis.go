package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis.UniversalClient
}

// HGetAllScan is a helper function to run HGETALL on a key and store the
// result in dest. dest must be an assingable address. The returned bool
// indicates if the key was found in the database.
func (c *Client) HGetAllScan(ctx context.Context, key string, dest any) (found bool, err error) {
	cmd := c.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	if len(cmd.Val()) == 0 {
		return false, nil
	}
	err = cmd.Scan(dest)
	if err != nil {
		return false, err
	}
	return true, nil
}
