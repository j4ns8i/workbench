package xredis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var _ redis.Hook = (*loggingHook)(nil)

type loggingHook struct {
	Logger *zerolog.Logger
}

// ProcessHook implements redis.Hook.
func (h *loggingHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)
		h.Logger.Info().Str("cmd", cmd.String()).Msg("executing command")
		return err
	}
}

// ProcessPipelineHook implements redis.Hook.
func (h *loggingHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		err := next(ctx, cmds)
		h.Logger.Info().Any("cmds", cmds).Msg("executing pipeline")
		return err
	}
}

func (h *loggingHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}
