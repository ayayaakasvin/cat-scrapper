package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	imagepool "github.com/ayayaakasvin/cat-photo-fetch/image-pool"
	apiserver "github.com/ayayaakasvin/cat-scrapper/internal/api-server"
	"github.com/ayayaakasvin/cat-scrapper/internal/api-server/libs/aliveapp"
	"github.com/ayayaakasvin/cat-scrapper/internal/config"
	saveengine "github.com/ayayaakasvin/cat-scrapper/internal/save-engine"
	"github.com/ayayaakasvin/cat-scrapper/pkg/logger"
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

	cfg := config.MustLoadConfig()
	log := logger.New(cfg.Logger)

	log.Info("config", "env", cfg.Logger.Env)

	sg, err := saveengine.NewSaveEngine(cfg.SavePath)
	pool, err := imagepool.NewCatImagePool()
	if err != nil {
		return fmt.Errorf("init error: %s", err)
	}

	gs := setupSupervisor(ctx, log)

	app := apiserver.NewApiServer(
		&cfg.HTTPServerConfig,
		&cfg.CorsConfig,
		log,
		sg,
		pool,
	)

	gs.Go("Server Status", aliveapp.LogAppStatus(time.Minute*3, log, ctx))
	gs.Go("Memory Stats", aliveapp.MemStat(time.Minute*3, log, ctx))
	gs.Go("http-server", app.Start)

	err = gs.Wait()
	if err != nil {
		return fmt.Errorf("gs wait error: %s", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Stop(shutdownCtx)

	return nil
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
