package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	apiserver "github.com/ayayaakasvin/cat-scrapper/internal/api-server"
)


func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := slog.Default()

	gs := setupSupervisor(ctx, log)

	app := apiserver.NewApiServer(
		&cfg.HTTPServer,
		&cfg.CorsConfig,
		cfg.GateawaySecret,
		log,
		repo,
		cache,
		jwtM,
	)

	gs.Go("http-server", app.Start)

	err = gs.Wait()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Stop(shutdownCtx)
	repo.Close()
	cache.Close()

	return err
}

func setupSupervisor(ctx context.Context, log *logrus.Logger) *goroutinesupervisor.GoRoutineSupervisor {
	gs := goroutinesupervisor.NewSupervisor(ctx)
	gs.WithHandler(func(e goroutinesupervisor.Event) {
		switch e.Type {
		case goroutinesupervisor.EventTaskStarted:
			log.Infof("Task %s started at %s", e.Task, e.Started.String())
		case goroutinesupervisor.EventTaskFinished:
			log.Infof("Task %s finished at %s", e.Task, e.Ended.String())
		case goroutinesupervisor.EventTaskFailed:
			log.Infof("Task %s failed at %s", e.Task, e.Ended.String())
		default:
		}
	})

	return gs
}