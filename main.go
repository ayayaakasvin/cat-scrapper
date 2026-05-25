package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiserver "github.com/ayayaakasvin/cat-scrapper/internal/api-server"
	"github.com/ayayaakasvin/cat-scrapper/internal/config"
	saveengine "github.com/ayayaakasvin/cat-scrapper/internal/save-engine"
	"github.com/ayayaakasvin/goroutinesupervisor"
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
	cfg := config.MustLoadConfig()

    sg, err := saveengine.NewSaveEngine(cfg.SavePath)

	gs := setupSupervisor(ctx, log)

	app := apiserver.NewApiServer(
		&cfg.HTTPServerConfig,
		log,
        sg,
	)

	gs.Go("http-server", app.Start)

	err = gs.Wait()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Stop(shutdownCtx)

	return err
}

func setupSupervisor(ctx context.Context, log *slog.Logger) *goroutinesupervisor.GoRoutineSupervisor {
	gs := goroutinesupervisor.NewSupervisor(ctx)
	gs.WithHandler(func(e goroutinesupervisor.Event) {
		switch e.Type {
		case goroutinesupervisor.EventTaskStarted:
			log.Info("Task started", "task", e.Task, "time", e.Started.String())
		case goroutinesupervisor.EventTaskFinished:
			log.Info("Task finished", "task", e.Task, "time", e.Ended.String())
		case goroutinesupervisor.EventTaskFailed:
			log.Info("Task failed", "task", e.Task, "time", e.Ended.String())
		default:
		}
	})

	return gs
}
