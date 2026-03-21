package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tgBookBot/internal/bot"
	"tgBookBot/internal/config"
	service "tgBookBot/internal/services"
	"tgBookBot/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	logger := setapLogger(cfg.Env)
	client, err := postgresql.NewClient(ctx, &cfg.DBConfig, logger)
	if err != nil {
		logger.Error("Failed to create database client", slog.Any("error", err))
		return
	}
	logger.Info("Database client created successfully")
	repo := postgresql.NewRepository(client, logger)
	logger.Info("Repository created successfully")
	service := service.New(logger, repo)
	logger.Info("Service created successfully")
	b := bot.New(&cfg.TGBotConfig, service, logger)
	logger.Info("Bot created successfully")
	//b.engine.Handle("/health", func(c telebot.Context) error {
	//return c.Send("OK")
	//})
	go b.Start()
	logger.Info("Bot started success")
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT)
	stopSignal := <-osSignals
	b.Stop()
	logger.Info("Bot graceful stopped", slog.String("signal", stopSignal.String()))
}
func setapLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
