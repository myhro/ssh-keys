package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/myhro/ssh-keys/config"
	"github.com/myhro/ssh-keys/syncer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = syncer.New(cfg).Run(ctx)
	if err != nil {
		slog.Error("sync failed", "error", err)
		os.Exit(1)
	}
}
