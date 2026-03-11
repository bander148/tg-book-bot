package service

import (
	"context"
	"fmt"
	"log/slog"
	model "tgBookBot/internal/domain/models"
	"tgBookBot/internal/storage"
	"time"
)

type Service struct {
	log     *slog.Logger
	storage storage.Storage
}

func New(log *slog.Logger, storage storage.Storage) *Service {
	return &Service{log: log, storage: storage}
}
func (s *Service) CreateUser(ctx context.Context, telegramID int64, username string) error {
	const op = "Service.CreateUser"
	log := s.log.With(slog.String("op", op), slog.Int64("telegramID", telegramID))
	log.Info("attempting to create new user")

	date := time.Now()
	if err := s.storage.CreateUser(ctx, telegramID, username, date); err != nil {
		log.Error("failed to create user", slog.Any("error", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	log.Info("user created")
	return nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	panic("implement me")
}

func (s *Service) CreateBook(ctx context.Context, book *model.Book) error {
	panic("implement me")
}

func (s *Service) GetUserBook(ctx context.Context, userID int64) ([]*model.Book, error) {
	panic("implement me")
}

func (s *Service) UpdateBookProgress(ctx context.Context, bookID int64, pagesRead int64) error {
	panic("implement me")
}
func (s *Service) MarkBookFinished(ctx context.Context, bookID int64) error {
	panic("implement me")
}
func (s *Service) DeleteBook(ctx context.Context, bookID int64) error {
	panic("implement me")
}
func (s *Service) GetBookByID(ctx context.Context, bookID int64) (*model.Book, error) {
	panic("implement me")
}
func (s *Service) AddReadSession(ctx context.Context, session *model.ReadingSession) error {
	panic("implement me")
}
func (s *Service) GetReadingSessionForPeriod(ctx context.Context, userID int64, from time.Time, to time.Time) ([]*model.ReadingSession, error) {
	panic("implement me")
}
