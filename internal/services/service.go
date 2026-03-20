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
func (s *Service) CreateUser(ctx context.Context, telegramID int64, username string) (int, error) {
	// TODO : request validation
	const op = "Service.CreateUser"
	l := s.log.With(slog.String("op", op), slog.Int64("telegramID", telegramID))
	l.Info("attempting to create new user")

	id, err := s.storage.CreateUser(ctx, telegramID, username)
	if err != nil {
		l.Error("failed to create user", slog.Any("error", err))
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	l.Info("user created", slog.Int("id", id))
	return id, nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	// TODO : request validation

	const op = "Service.GetUserByTelegramID"
	l := s.log.With(slog.String("op", op), slog.Int64("telegramID", telegramID))
	l.Info("attempting to get user info")

	userData, err := s.storage.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		l.Error("Failed to get user info", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	// think about where and how to process empty data from the database
	return userData, nil
}

func (s *Service) CreateBook(ctx context.Context, book *model.Book) error {
	panic("implement me")
}

func (s *Service) GetUserBooks(ctx context.Context, userID int64) ([]model.Book, error) {
	const op = "Service.GetUserBook"
	l := s.log.With(slog.String("op", op), slog.Int64("userID", userID))
	l.Info("attempting to get user books")

	booksData, err := s.storage.GetUserBooks(ctx, userID)
	if err != nil {
		l.Error("Failed to get user books", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	// think about where and how to process empty data from the database
	return booksData, nil
}

func (s *Service) UpdateBookProgress(ctx context.Context, bookID int64, pagesRead int64) error {
	panic("implement me")
}
func (s *Service) MarkBookFinished(ctx context.Context, bookID int64, date time.Time) error {
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
