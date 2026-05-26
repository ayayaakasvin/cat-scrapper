package aliveapp

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

func LogAppStatus(trate time.Duration, log *slog.Logger, ctx context.Context) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(trate)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Info("Server is alive...")
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MemStat(trate time.Duration, log *slog.Logger, ctx context.Context) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(trate)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				log.Info("Allocated mem","MiB", m.Alloc/1024/1024)
				time.Sleep(1 * time.Second)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
