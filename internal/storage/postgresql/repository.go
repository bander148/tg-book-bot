package postgresql

import (
	"context"
	"log/slog"
	model "tgBookBot/internal/domain/models"
	"tgBookBot/internal/storage"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type pgRepository struct {
	client Client
	log    *slog.Logger
}

func NewRepository(client Client, log *slog.Logger) storage.Storage {
	return &pgRepository{client: client, log: log}
}

// TODO : create specific errors for repository and handle them in service layer, for example ErrUserNotFound, ErrBookNotFound etc.

// TODO : use dto for operations - userDTO
func (p *pgRepository) CreateUser(ctx context.Context, telegramID int64, username string) (int, error) {
	const op = "pdRepository.CreateUser"
	var id int
	l := p.log.With(slog.String("op", op), slog.Int64("tg_id", telegramID))
	q := `INSERT INTO users (telegram_id,username) VALUES ($1,$2) RETURNING id`
	if err := p.client.QueryRow(ctx, q, telegramID, username).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't save user", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return 0, pgErr
		}

		return 0, err
	}
	l.Info("user created", slog.Int("id", id))
	return id, nil
}
func (p *pgRepository) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	const op = "pdRepository.GetUserByTelegramID"
	l := p.log.With(slog.String("op", op), slog.Int64("telegramID", telegramID))
	q := `SELECT id , telegram_id , username , created_at FROM users WHERE telegram_id = $1`
	var user model.User
	if err := p.client.QueryRow(ctx, q, telegramID).Scan(&user.ID, &user.TelegramID, &user.Username, &user.CreatedAt); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't find user", slog.Any("err", pgErr.Message), slog.Any("deatil", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return nil, pgErr
		}
		return nil, err
	}
	return &user, nil

}
func (p *pgRepository) CreateBook(ctx context.Context, book *model.Book) error {
	panic("implement me ")
}
func (p *pgRepository) GetUserBooks(ctx context.Context, userID int64) ([]model.Book, error) {
	const op = "pdRepository.GetUserBooks"
	l := p.log.With(slog.String("op", op), slog.Int64("tg_id", userID))

	q := `SELECT id,user_id , pages , description , author , title , start_date , end_date, pages_read, created_at, updated_at FROM books WHERE user_id = $1` // TODO : make sql request , указывать поля а не *
	rows, err := p.client.Query(ctx, q, userID)
	if err != nil {
		l.Error("can't scan row", slog.Any("err", err))
		return []model.Book{}, err
	}

	books := make([]model.Book, 0)

	for rows.Next() {
		var book model.Book

		err := rows.Scan(&book.ID, &book.UserID, &book.Pages, &book.Description, &book.Author, &book.Title, &book.StartDate, &book.EndDate, &book.PagesRead, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			l.Error("can't scan row", slog.Any("err", err))
			return []model.Book{}, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		l.Error("rows error", slog.Any("err", err))
		return []model.Book{}, err
	}
	return books, nil

}
func (p *pgRepository) UpdateBookProgress(ctx context.Context, bookID int64, pagesRead int64) error {
	panic("implement me ")
}
func (p *pgRepository) MarkBookFinished(ctx context.Context, bookID int64, date time.Time) error {
	panic("implement me ")
}
func (p *pgRepository) DeleteBook(ctx context.Context, bookID int64) error {
	const op = "pdRepository.DeleteBook"
	l := p.log.With(slog.String("op", op), slog.Int64("book_id", bookID))

	q := "DELETE FROM books WHERE id = $1"
	if _, err := p.client.Exec(ctx, q, bookID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't delete user", slog.Any("err", pgErr.Message), slog.Any("deatil", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return pgErr
		}
		return err
	}
	return nil
}
func (p *pgRepository) GetBookByID(ctx context.Context, bookID int64) (*model.Book, error) {
	panic("implement me ")
}
func (p *pgRepository) AddReadingSession(ctx context.Context, session *model.ReadingSession) error {
	panic("implement me ")
}
func (p *pgRepository) GetReadingSessionsForPeriod(
	ctx context.Context,
	userID int64,
	from time.Time,
	to time.Time,
) ([]*model.ReadingSession, error) {
	panic("implement me ")
}
