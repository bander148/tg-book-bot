package service

import (
	"context"
	"fmt"
	"log/slog"
	model "tgBookBot/internal/domain/models"
	"tgBookBot/internal/dto"
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
func (s *Service) CreateUser(ctx context.Context, req *dto.UserCreateRequest) error {
	const op = "Service.CreateUser"
	l := s.log.With(slog.String("op", op), slog.Int64("telegram_id", req.TelegramID))

	user := &model.User{
		TelegramID: req.TelegramID,
		Username:   req.Username,
	}

	id, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		l.Error("failed to create user", slog.Any("error", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	l.Info("user created", slog.Int("id", id))
	return nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	// TODO : request validation

	const op = "Service.GetUserByTelegramID"
	l := s.log.With(slog.String("op", op), slog.Int64("telegram_id", telegramID))

	userData, err := s.storage.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		l.Error("Failed to get user info", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	// think about where and how to process empty data from the database
	l.Info("user data getted", slog.Any("username", userData.Username))
	return userData, nil
}

func (s *Service) CreateBook(ctx context.Context, book *dto.BookCreateRequest) error {
	const op = "Service.CreateBook"
	// TODO : request validation
	l := s.log.With(slog.String("op", op), slog.Int64("user_id", book.TelegramID))
	userID, err := s.storage.GetUserIDByTelegramID(ctx, book.TelegramID)
	if err != nil {
		l.Error("Failed to get user ID", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	bookModel := &model.Book{
		Pages:       book.Pages,
		Description: book.Description,
		Author:      book.Author,
		Title:       book.Title,
		UserID:      userID,
		PagesRead:   book.PagesRead,
		StartDate:   nil,
		EndDate:     nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := s.storage.CreateBook(ctx, bookModel)
	if err != nil {
		l.Error("failed to create book", slog.Any("error", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	l.Info("book created", slog.Int("id", id))
	return nil

}

func (s *Service) GetUserBooks(ctx context.Context, userID int64) ([]model.Book, error) {
	const op = "Service.GetUserBook"
	l := s.log.With(slog.String("op", op), slog.Int64("user_id", userID))
	// TODO : request validation

	booksData, err := s.storage.GetUserBooks(ctx, userID)
	if err != nil {
		l.Error("Failed to get user books", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	// think about where and how to process empty data from the database
	l.Info("user books getted", slog.Int("books_count", len(booksData)))
	return booksData, nil
}

func (s *Service) UpdateBookProgress(ctx context.Context, book *dto.BookUpdateProgressRequest) (*dto.BookResponse, error) {
	const op = "Service.UpdateBookProgress"
	l := s.log.With(slog.String("op", op), slog.Int64("book_id", book.BookID), slog.Int64("telegram_id", book.TelegramID))
	if err := s.CheckAccessToBook(ctx, book.BookID, book.TelegramID); err != nil {
		l.Error("Access to book denied", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	// TODO : request validation
	bookModel := &model.Book{
		ID:        book.BookID,
		PagesRead: book.PagesRead,
		UpdatedAt: time.Now(),
	}
	bookResponse, err := s.storage.UpdateBookProgress(ctx, bookModel)
	if err != nil {
		l.Error("Failed to update book progress", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	l.Info("book progress updated", slog.Int64("book_id", book.BookID))

	return &dto.BookResponse{
		ID:          bookResponse.ID,
		Pages:       bookResponse.Pages,
		Description: bookResponse.Description,
		Author:      bookResponse.Author,
		Title:       bookResponse.Title,
		PagesRead:   bookResponse.PagesRead,
		StartDate:   bookResponse.StartDate,
		EndDate:     bookResponse.EndDate,
		CreatedAt:   bookResponse.CreatedAt,
		UpdatedAt:   bookResponse.UpdatedAt,
	}, nil

}
func (s *Service) MarkBookFinished(ctx context.Context, bookDTO *dto.BookMarkFinishedRequest) (*dto.BookResponse, error) {
	const op = "Service.MarkBookFinished"
	l := s.log.With(slog.String("op", op), slog.Int64("book_id", bookDTO.BookID), slog.Int64("telegram_id", bookDTO.TelegramID))
	// TODO : request validation
	if err := s.CheckAccessToBook(ctx, bookDTO.BookID, bookDTO.TelegramID); err != nil {
		l.Error("Access to book denied", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	bookModel := &model.Book{
		ID:      bookDTO.BookID,
		EndDate: &bookDTO.EndDate,
	}
	bookResponse, err := s.storage.MarkBookFinished(ctx, bookModel)
	if err != nil {
		l.Error("Failed to mark book as finished", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &dto.BookResponse{
		ID:          bookResponse.ID,
		Pages:       bookResponse.Pages,
		Description: bookResponse.Description,
		Author:      bookResponse.Author,
		Title:       bookResponse.Title,
		PagesRead:   bookResponse.PagesRead,
		StartDate:   bookResponse.StartDate,
		EndDate:     bookResponse.EndDate,
		CreatedAt:   bookResponse.CreatedAt,
		UpdatedAt:   bookResponse.UpdatedAt,
	}, nil

}
func (s *Service) DeleteBook(ctx context.Context, bookDTO *dto.BookDeleteRequest) error {
	const op = "Service.DeleteBook"
	l := s.log.With(slog.String("op", op), slog.Int64("book_id", bookDTO.BookID), slog.Int64("telegram_id", bookDTO.TelegramID))
	// TODO : request validation
	if err := s.CheckAccessToBook(ctx, bookDTO.BookID, bookDTO.TelegramID); err != nil {
		l.Error("Access to book denied", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	if err := s.storage.DeleteBook(ctx, bookDTO.BookID); err != nil {
		l.Error("Failed to delete book", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	l.Info("book deleted", slog.Int64("book_id", bookDTO.BookID))
	return nil
}
func (s *Service) GetBookByID(ctx context.Context, bookID int64) (*model.Book, error) {
	const op = "Service.GetBookByID"
	l := s.log.With(slog.String("op", op), slog.Int64("book_id", bookID))
	// TODO : request validation

	bookData, err := s.storage.GetBookByID(ctx, bookID)
	if err != nil {
		l.Error("Failed to get book info", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return bookData, nil
}
func (s *Service) AddReadSession(ctx context.Context, sessionDTO *dto.ReadingSessionRequest) error {
	const op = "Service.AddReadSession"
	l := s.log.With(slog.String("op", op), slog.Int64("user_id", sessionDTO.TelegramID), slog.Int64("book_id", sessionDTO.BookID))
	// TODO : request validation
	userID, err := s.storage.GetUserIDByTelegramID(ctx, sessionDTO.TelegramID)
	if err != nil {
		l.Error("Failed to get user ID", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	sessionModel := &model.ReadingSession{
		UserID:   userID,
		BookID:   sessionDTO.BookID,
		Pages:    sessionDTO.Pages,
		Date:     sessionDTO.Date,
		Duration: sessionDTO.Duration,
	}
	if err := s.storage.AddReadingSession(ctx, sessionModel); err != nil {
		l.Error("Failed to add reading session", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil

}
func (s *Service) GetReadingSessionForPeriod(ctx context.Context, sessionDTO *dto.ReadingSessionsForPeriodRequest) ([]model.ReadingSession, error) {
	const op = "Service.GetReadingSessionForPeriod"
	l := s.log.With(slog.String("op", op), slog.Int64("user_id", sessionDTO.TelegramID), slog.Time("from", sessionDTO.From), slog.Time("to", sessionDTO.To))
	// TODO : request validation
	userID, err := s.storage.GetUserIDByTelegramID(ctx, sessionDTO.TelegramID)
	if err != nil {
		l.Error("Failed to get user ID", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	sessions, err := s.storage.GetReadingSessionsForPeriod(ctx, sessionDTO.From, sessionDTO.To, userID)
	if err != nil {
		l.Error("Failed to get reading sessions for period", slog.Any("err", err))
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return sessions, nil

}

// Think about where to place this method , maybe refactor it in future
func (s *Service) CheckAccessToBook(ctx context.Context, bookID int64, telegramID int64) error {
	const op = "Service.CheckAccessToBook"
	l := s.log.With(slog.String("op", op), slog.Int64("book_id", bookID), slog.Int64("telegram_id", telegramID))

	userID, err := s.storage.GetUserIDByTelegramID(ctx, telegramID)
	if err != nil {
		l.Error("Failed to get user ID", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}

	bookStorage, err := s.storage.GetBookByID(ctx, bookID)
	if err != nil {
		l.Error("Failed to get book ID", slog.Any("err", err))
		return fmt.Errorf("%s:%w", op, err)
	}

	if bookStorage.UserID != userID {
		l.Error("User ID from telegram does not match book owner ID", slog.Int64("book_user_id", bookStorage.UserID), slog.Int64("user_id", userID))
		return fmt.Errorf("%s: %w", op, fmt.Errorf("user ID from telegram does not match book owner ID"))
	}
	return nil
}
