package bot

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"tgBookBot/internal/dto"

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
		req := &dto.UserCreateRequest{TelegramID: telegram_id, Username: username}
		err := b.service.CreateUser(ctx, req)
		if err != nil {
			return c.Send("Произошла ошибка при регистрации. Пожалуйста, попробуйте снова позже.")
		}
	}
	b.log.Info("User already exists", slog.Int64("telegram_id", telegram_id), slog.String("username", username)) // delelte in future, just for health check
	return c.Send("Привет! Я бот для отслеживания прогресса чтения книг. Я помогу тебе следить за тем, сколько страниц ты прочитал и когда ты закончил книгу. Чтобы начать, просто добавь книгу с помощью команды /addbook.")
}
func (b *Bot) handleAddBook(c telebot.Context) error {
	// TODO : add implementation of context with timeout and cancellation
	ctx := context.Background()
	telegramID := c.Sender().ID
	text := strings.TrimSpace(c.Message().Payload)
	if text == "" {
		return c.Send("Использование: /addbook Название | Автор | Количество страниц")
	}
	parts := strings.Split(text, "|")
	if len(parts) < 3 {
		return c.Send("Неверный формат.\n Пример:\n/addbook Война и мир | Лев Толстой | 500")
	}
	title := strings.TrimSpace(parts[0])
	author := strings.TrimSpace(parts[1])
	pagesStr := strings.TrimSpace(parts[2])

	pages, err := strconv.ParseInt(pagesStr, 10, 64)
	if err != nil || pages <= 0 {
		return c.Send("Количество страниц должно быть положительным числом.")
	}
	dto := &dto.BookCreateRequest{
		TelegramID:  telegramID,
		Title:       title,
		Author:      author,
		Pages:       pages,
		Description: "", // delete from dto, because it's not required field, maybe implement this functionality in future
		PagesRead:   0,
	}
	err = b.service.CreateBook(ctx, dto)
	if err != nil {
		b.log.Error("Failed to create book", slog.Any("err", err))
		return c.Send("Произошла ошибка при добавлении книги. Пожалуйста, попробуйте снова позже.")
	}
	return c.Send("Книга успешно добавлена!")
}
