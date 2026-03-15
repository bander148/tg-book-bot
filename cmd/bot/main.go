package main

import (
	"log/slog"
	"os"
	"tgBookBot/internal/config"

	"gopkg.in/telebot.v4"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := setapLogger(cfg.Env)
	b, err := telebot.NewBot(telebot.Settings{Token: cfg.TGBotConfig.APIToken, Poller: &telebot.LongPoller{Timeout: cfg.TGBotConfig.PollTime}})
	if err != nil {
		logger.Error("can't start bot", slog.Any("err", err))
		panic("Failed to start bot")
	}
	logger.Info("Bot started success")
	b.Handle("/health", func(c telebot.Context) error {
		return c.Send("OK")
	})
	b.Start()
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
