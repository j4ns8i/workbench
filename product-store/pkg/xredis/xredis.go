package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis.UniversalClient
}

// HGetAllScan is a helper function to run HGETALL on a key and store the
// result in dest. dest must be an assignable address. The returned bool
// indicates if the key was found in the database.
func HGetAllScan(ctx context.Context, client redis.HashCmdable, key string, dest any) (bool, error) {
	cmd := client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	if len(cmd.Val()) == 0 {
		return false, nil
	}
	err := cmd.Scan(dest)
	if err != nil {
		return false, err
	}
	return true, nil
}
