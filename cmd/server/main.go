package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ds9labs.com/space-launch-server/internal/config"
	"ds9labs.com/space-launch-server/internal/handler"
	"ds9labs.com/space-launch-server/internal/httpserver"
	"ds9labs.com/space-launch-server/internal/notifications"
	"ds9labs.com/space-launch-server/internal/provider/spacedevs"
	"ds9labs.com/space-launch-server/internal/service"
	"ds9labs.com/space-launch-server/internal/storage"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store, err := storage.NewValkeyStore(cfg.Valkey)
	if err != nil {
		logger.Error("failed to initialize valkey store", "error", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := store.Close(); cerr != nil {
			logger.Warn("failed to close valkey store", "error", cerr)
		}
	}()

	spaceDevsClient := spacedevs.NewClient(cfg.SpaceDevs, store)
	upcomingService := service.NewSpaceDevsUpcomingService(spaceDevsClient, store, cfg.SpaceDevs)
	lcdSpaceService := service.NewLCDSpaceService(upcomingService, store, service.LCDSpaceConfig{}, nil)
	registry := notifications.NewRegistry(store)

	var sender notifications.Sender = notifications.DisabledSender{}
	if cfg.Notifications.Enabled {
		firebaseSender, senderErr := notifications.NewFirebaseSender(ctx, cfg.Notifications.FirebaseProjectID, cfg.Notifications.CredentialsFile)
		if senderErr != nil {
			logger.Error("FCM is enabled but sender initialization failed", "error", senderErr)
			os.Exit(1)
		}
		sender = firebaseSender
		go notifications.NewMonitor(upcomingService, registry, sender, logger, cfg.SpaceDevs.LaunchesLimit).Run(ctx, cfg.Notifications.PollInterval)
		logger.Info("FCM notification monitor enabled", "poll_interval", cfg.Notifications.PollInterval)
	}

	api := handler.NewAPI(logger, store, upcomingService, registry, lcdSpaceService)
	router := httpserver.NewRouter(api)
	srv := httpserver.New(cfg.HTTP, router)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting server", "addr", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server terminated", "error", err)
			os.Exit(1)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server shutdown complete")
}
