package storage

import (
	"context"
	model "tgBookBot/internal/domain/models"
	"time"
)

type Storage interface {

	// User operations
	CreateUser(ctx context.Context, telegramID int64, username string, createdAt time.Time) error
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)

	// Operations with books
	CreateBook(ctx context.Context, book *model.Book) error
	GetUserBook(ctx context.Context, userID int64) ([]*model.Book, error)
	UpdateBookProgress(ctx context.Context, bookID int64, pagesRead int64) error
	MarkBookFinished(ctx context.Context, bookID int64, date time.Time) error
	DeleteBook(ctx context.Context, bookID int64) error
	GetBookByID(ctx context.Context, bookID int64) (*model.Book, error)

	// Operations with ReadingSession
	AddReadingSession(ctx context.Context, session *model.ReadingSession) error
	GetReadingSessionsForPeriod(ctx context.Context, userID int64, from time.Time, to time.Time) ([]*model.ReadingSession, error)
}
