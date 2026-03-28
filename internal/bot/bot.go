package bot

import (
	"log/slog"
	"tgBookBot/internal/config"
	service "tgBookBot/internal/services"
	"time"

	"gopkg.in/telebot.v4"
)

// Bot - основной тип, представляющий бота. Содержит поля для движка бота (engine), сервиса (service) и логгера (log).
type Bot struct {
	engine  *telebot.Bot
	service *service.Service
	log     *slog.Logger
}

// New - конструктор для создания нового экземпляра бота. Принимает конфигурацию, сервис и логгер, и возвращает указатель на новый экземпляр Bot
func New(cfg *config.TelegramBotConfig, service *service.Service, log *slog.Logger) *Bot {

	engine := CreateEngine(cfg.APIToken, cfg.PollTime)

	b := &Bot{
		engine:  engine,
		service: service,
		log:     log,
	}
	log.Info("Bot created success", slog.Any("poll_time", cfg.PollTime))
	b.registerHandlers()
	return b
}

// Start - запускает бота и начинает обрабатывать входящие сообщения
func (b *Bot) Start() {
	b.engine.Start()
}

// Stop - останавливает бота и прекращает обработку входящих сообщений , используется для graceful shutdown
func (b *Bot) Stop() {
	b.engine.Stop()
}

// CreateEngine - функция для создания нового экземпляра telebot.Bot с заданным токеном и временем опроса. Если при создании бота возникает ошибка, функция паникует
func CreateEngine(token string, poll time.Duration) *telebot.Bot {
	engine, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: poll},
	})
	if err != nil {
		panic("Failed to start bot: " + err.Error())
	}
	return engine

}

// registerHandlers - метод для регистрации обработчиков команд и сообщений бота
func (b *Bot) registerHandlers() {
	// TODO: register handlers
	b.engine.Handle("/start", b.handleStart)
	b.engine.Handle("/addbook", b.handleAddBook)
	// TODO : add handlers for other commands - /addbook, /mybooks, /deletebook etc.
	// TODO : add handlers for callback buttons
	//panic("implement me ")
}
