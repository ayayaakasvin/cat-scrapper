package wplog

import (
	"context"
	"log/slog"

	"github.com/ayayaakasvin/wpn"
)

func WorkerPoolResultLogger(chanResult <-chan *wpn.Result, log *slog.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case result, ok := <-chanResult:
				if !ok {
					return ctx.Err()
				}

				if result == nil {
					log.Warn("received nil worker pool result")
					continue
				}

				attrs := []slog.Attr{
					slog.String("job_id", result.JobID),
					slog.Int("worker_id", result.WorkerID),
					slog.Int("attempts", result.Attempts),
					slog.String("output", result.Output.String()),
					slog.Duration("duration", result.TimeConsumed),
					slog.Time("started_at", result.StartedAt),
					slog.Time("finished_at", result.FinishedAt),
				}

				if result.Error != nil {
					log.LogAttrs(ctx, slog.LevelError, "job failed", append(attrs, slog.Any("error", result.Error))...)
					continue
				}

				log.LogAttrs(ctx, slog.LevelInfo, "job completed", attrs...)
			}
		}
	}
}
