package storage

import (
	"context"
	model "tgBookBot/internal/domain/models"
	"time"
)

// TODO : differentiate between user and book operations in storage interface, maybe create separate interfaces for them
type Storage interface {

	// User operations
	CreateUser(ctx context.Context, user *model.User) (int, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)

	// Operations with books
	CreateBook(ctx context.Context, book *model.Book) (int, error)
	GetUserBooks(ctx context.Context, userID int64) ([]model.Book, error)
	UpdateBookProgress(ctx context.Context, book *model.Book) (*model.Book, error)
	MarkBookFinished(ctx context.Context, book *model.Book) (*model.Book, error)
	DeleteBook(ctx context.Context, bookID int64) error
	GetBookByID(ctx context.Context, bookID int64) (*model.Book, error)

	// Operations with ReadingSession
	AddReadingSession(ctx context.Context, session *model.ReadingSession) error
	GetReadingSessionsForPeriod(ctx context.Context, from time.Time, to time.Time, userID int64) ([]model.ReadingSession, error)
}
