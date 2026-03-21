package bot

import (
	"context"
	"log/slog"

	"gopkg.in/telebot.v4"
)

func (b *Bot) handleStart(c telebot.Context) error {
	username := c.Sender().Username
	if username == "" {
		username = c.Sender().FirstName
	}
	telegram_id := c.Sender().ID
	ctx := context.Background() // TODO : add implementation of context with timeout and cancellation
	ok, err := b.service.GetUserByTelegramID(ctx, telegram_id)
	if err != nil {
		b.log.Info("Failed to get user info", slog.Any("err", err))
	}
	if ok == nil {
		b.log.Info("Creating new user", slog.Int64("telegram_id", telegram_id), slog.String("username", username))
		id, err := b.service.CreateUser(ctx, telegram_id, username)
		if err != nil {
			return c.Send("Произошла ошибка. Пожалуйста, попробуйте снова позже.")
		}
		if id == 0 {
			return c.Send("Произошла ошибка. Пожалуйста, попробуйте снова позже.")
		}
	}
	return c.Send("Привет! Я бот для отслеживания прогресса чтения книг. Я помогу тебе следить за тем, сколько страниц ты прочитал и когда ты закончил книгу. Чтобы начать, просто добавь книгу с помощью команды /addbook.")
}
