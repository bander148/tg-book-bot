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

func (p *pgRepository) CreateUser(ctx context.Context, user *model.User) (int, error) {
	const op = "pdRepository.CreateUser"
	var id int
	l := p.log.With(slog.String("op", op), slog.Int64("tg_id", user.TelegramID))
	q := `INSERT INTO users (telegram_id,username) VALUES ($1,$2) RETURNING id`
	if err := p.client.QueryRow(ctx, q, user.TelegramID, user.Username).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't save user", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return 0, pgErr
		}

		return 0, err
	}
	user.ID = int64(id)
	l.Info("user created", slog.Int("id", id))
	return id, nil
}
func (p *pgRepository) GetUserIDByTelegramID(ctx context.Context, telegramID int64) (int64, error) {
	panic("implement me")
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
func (p *pgRepository) CreateBook(ctx context.Context, book *model.Book) (int, error) {
	const op = "pgRepository.CreateBook"
	l := p.log.With(slog.String("op", op), slog.Int64("user_id", book.UserID), slog.String("title", book.Title))
	q := `INSERT INTO books(pages,description,author,title,user_id,start_date,end_date,pages_read) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	if err := p.client.QueryRow(ctx, q, book.Pages, book.Description, book.Author, book.Title, book.UserID, book.StartDate, book.EndDate, book.PagesRead).Scan(&book.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't save book", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return 0, pgErr
		}
		return 0, err
	}
	return int(book.ID), nil

}
func (p *pgRepository) GetUserBooks(ctx context.Context, userID int64) ([]model.Book, error) {
	const op = "pdRepository.GetUserBooks"
	l := p.log.With(slog.String("op", op), slog.Int64("tg_id", userID))

	q := `SELECT id,user_id , pages , description , author , title , start_date , end_date, pages_read, created_at, updated_at FROM books WHERE user_id = $1`
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
func (p *pgRepository) UpdateBookProgress(ctx context.Context, book *model.Book) (*model.Book, error) {
	const op = "pdRepository.UpdateBookProgress"

	l := p.log.With(slog.String("op", op), slog.Int64("book_id", book.ID), slog.Int64("user_id", book.UserID))
	q := `UPDATE books SET pages_read = $1, updated_at = $2 WHERE id = $3 RETURNING id,user_id , pages , description , author , title , start_date , end_date, pages_read, created_at, updated_at`

	var updateBook model.Book
	if err := p.client.QueryRow(ctx, q, book.PagesRead, time.Now(), book.ID).Scan(&updateBook.ID, &updateBook.UserID, &updateBook.Pages, &updateBook.Description, &updateBook.Author, &updateBook.Title, &updateBook.StartDate, &updateBook.EndDate, &updateBook.PagesRead, &updateBook.CreatedAt, &updateBook.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't update book progress", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return nil, pgErr
		}
		return nil, err
	}
	return &updateBook, nil
}
func (p *pgRepository) MarkBookFinished(ctx context.Context, book *model.Book) (*model.Book, error) {
	const op = "pdRepository.MarkBookFinished"
	l := p.log.With(slog.String("op", op), slog.Int64("book_id", book.ID))
	q := `UPDATE books SET end_date = $1, updated_at = $2 WHERE id = $3`
	if err := p.client.QueryRow(ctx, q, book.EndDate, time.Now(), book.ID).Scan(&book.ID, &book.UserID, &book.Pages, &book.Description, &book.Author, &book.Title, &book.StartDate, &book.EndDate, &book.PagesRead, &book.CreatedAt, &book.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't mark book as finished", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return nil, pgErr
		}
		return nil, err
	}
	return book, nil
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
	const op = "pdRepository.GetBookByID"
	l := p.log.With(slog.String("op", op), slog.Int64("book_id", bookID))
	q := `SELECT id,user_id , pages , description , author , title , start_date , end_date, pages_read, created_at, updated_at FROM books WHERE id = $1`
	var book model.Book
	if err := p.client.QueryRow(ctx, q, bookID).Scan(&book.ID, &book.UserID, &book.Pages, &book.Description, &book.Author, &book.Title, &book.StartDate, &book.EndDate, &book.PagesRead, &book.CreatedAt, &book.UpdatedAt); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't get book by ID", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return nil, pgErr
		}
		return nil, err
	}
	return &book, nil

}
func (p *pgRepository) AddReadingSession(ctx context.Context, session *model.ReadingSession) error {
	const op = "pdRepository.AddReadingSession"
	l := p.log.With(slog.String("op", op), slog.Int64("user_id", session.UserID), slog.Int64("book_id", session.BookID), slog.Any("date", session.Date), slog.Any("duration", session.Duration), slog.Int64("pages", session.Pages))
	q := `INSERT INTO reading_sessions (book_id, user_id, duration, date, pages) VALUES ($1, $2, $3, $4, $5)`
	if _, err := p.client.Exec(ctx, q, session.BookID, session.UserID, session.Duration, session.Date, session.Pages); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't add reading session", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return pgErr
		}
		return err
	}
	return nil
}
func (p *pgRepository) GetReadingSessionsForPeriod(ctx context.Context, from time.Time, to time.Time, userID int64) ([]model.ReadingSession, error) {
	const op = "pdRepository.GetReadingSessionsForPeriod"
	l := p.log.With(slog.String("op", op), slog.Int64("user_id", userID))
	q := `SELECT id, book_id, user_id, duration, date, pages, created_at, updated_at FROM reading_sessions WHERE user_id = $1 AND date >= $2 AND date <= $3`
	rows, err := p.client.Query(ctx, q, userID, from, to)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			l.Error("can't get reading sessions for period", slog.Any("err", pgErr.Message), slog.Any("detail", pgErr.Detail), slog.Any("where", pgErr.Where), slog.Any("code", pgErr.Code), slog.Any("sqlstate", pgErr.SQLState()))
			return []model.ReadingSession{}, pgErr
		}
		return []model.ReadingSession{}, err
	}

	sessions := make([]model.ReadingSession, 0)
	for rows.Next() {
		var session model.ReadingSession

		if err := rows.Scan(&session.ID, &session.BookID, &session.UserID, &session.Duration, &session.Date, &session.Pages, &session.CreatedAt, &session.UpdatedAt); err != nil {
			l.Error("can't scan row", slog.Any("err", err))
			return []model.ReadingSession{}, err
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		l.Error("rows error", slog.Any("err", err))
		return []model.ReadingSession{}, err
	}
	return sessions, nil

}

// TODO : add new methods for operations with reading sessions and books.
