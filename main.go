package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	imagepool "github.com/ayayaakasvin/cat-photo-fetch/image-pool"
	apiserver "github.com/ayayaakasvin/cat-scrapper/internal/api-server"
	"github.com/ayayaakasvin/cat-scrapper/internal/api-server/libs/appstat"
	"github.com/ayayaakasvin/cat-scrapper/internal/config"
	fsengine "github.com/ayayaakasvin/cat-scrapper/internal/fs_engine"
	"github.com/ayayaakasvin/cat-scrapper/internal/repository/sqlite"
	"github.com/ayayaakasvin/cat-scrapper/internal/wplog"
	"github.com/ayayaakasvin/wpn"

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

	sg, err := fsengine.NewFSE(cfg.SavePath)
	if err != nil {
		return fmt.Errorf("init save engine error: %w", err)
	}

	pool, err := imagepool.NewCatImagePool()
	if err != nil {
		return fmt.Errorf("init image pool error: %w", err)
	}
	wp, err := wpn.NewWorkerPool(runtime.NumCPU() * 2)
	if err != nil {
		return fmt.Errorf("worker pool init error: %s", err)
	}

	repo, err := sqlite.NewSqliteConnection(filepath.Join(cfg.SavePath, cfg.SqLiteConfig))
	if err != nil {
		return fmt.Errorf("init repository error: %w", err)
	}

	gs := setupSupervisor(ctx, log)

	app := apiserver.NewApiServer(
		&cfg.HTTPServerConfig,
		&cfg.CorsConfig,
		log,
		sg,
		repo,
		pool,
		wp,
	)

	gs.Go("Server Status", appstat.LogAppStatus(time.Minute*3, log, ctx))
	gs.Go("Memory Stats", appstat.MemStat(time.Minute*3, log, ctx))
	gs.Go("Worker Pool", wp.Start)
	gs.Go("Worker Pool Logger", wplog.WorkerPoolResultLogger(wp.Results(), log))
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
