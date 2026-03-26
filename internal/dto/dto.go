package dto

import "time"

// TODO : maybe split dto into separate files for user and book operations, for example user_dto.go, book_dto.go etc.

// BookCreateRequest for creating new book
type BookCreateRequest struct {
	Pages       int64  `json:"pages" validate:"required,gt=0"`
	Description string `json:"description" validate:"omitempty,max=2000"`
	Author      string `json:"author" validate:"required,max=100"`
	Title       string `json:"title" validate:"required,max=100"`
	TelegramID  int64  `json:"telegram_id" validate:"required,gt=0"`
	PagesRead   int64  `json:"pages_read" validate:"omitempty,gt=0"`
}

// BookUpdateRequest for updating book info
type UpdateBookRequest struct {
	BookID      int64  `json:"book_id" validate:"required,gt=0"`
	Pages       int64  `json:"pages" validate:"omitempty,gt=0"`
	Description string `json:"description" validate:"omitempty,max=2000"`
	Author      string `json:"author" validate:"omitempty,max=100"`
	Title       string `json:"title" validate:"omitempty,max=100"`
	TelegramID  int64  `json:"telegram_id" validate:"required,gt=0"`
}

// BookUpdateProgressRequest for updating book progress
type BookUpdateProgressRequest struct {
	BookID     int64 `json:"book_id" validate:"required,gt=0"`
	PagesRead  int64 `json:"pages_read" validate:"required,gt=0"`
	TelegramID int64 `json:"telegram_id" validate:"required,gt=0"`
}

// BookFinishRequest for marking book as finished
type BookMarkFinishedRequest struct {
	BookID     int64     `json:"book_id" validate:"required,gt=0"`
	EndDate    time.Time `json:"date" validate:"required"`
	TelegramID int64     `json:"telegram_id" validate:"required,gt=0"`
}

// BookDeleteRequest for deleting book
type BookDeleteRequest struct {
	BookID     int64 `json:"book_id" validate:"required,gt=0"`
	TelegramID int64 `json:"telegram_id" validate:"required,gt=0"`
}

// BookResponse for returning book info
type BookResponse struct {
	ID          int64      `json:"id"`
	Pages       int64      `json:"pages"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	PagesRead   int64      `json:"pages_read"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserCreateRequest for creating new user
type UserCreateRequest struct {
	TelegramID int64  `json:"telegram_id" validate:"required,gt=0"`
	Username   string `json:"username" validate:"required,max=100"`
}

// ReadingSessionRequest for creating new reading session
type ReadingSessionRequest struct {
	Pages      int64         `json:"pages" validate:"required,gt=0"`
	BookID     int64         `json:"book_id" validate:"required,gt=0"`
	TelegramID int64         `json:"telegram_id" validate:"required,gt=0"`
	Date       time.Time     `json:"date" validate:"required"`
	Duration   time.Duration `json:"duration" validate:"required,gt=0"`
}

// ReadingSessionResponse for returning reading session info
type ReadingSessionResponse struct {
	BookID   int64         `json:"book_id"`
	Duration time.Duration `json:"duration"`
	Date     time.Time     `json:"date"`
	Pages    int64         `json:"pages"`
}

// ReadingSessionsForPeriodRequest for getting reading sessions for period
type ReadingSessionsForPeriodRequest struct {
	TelegramID int64     `json:"telegram_id" validate:"required,gt=0"`
	From       time.Time `json:"from" validate:"required"`
	To         time.Time `json:"to" validate:"required,gtfield=From"`
}

// ReadingSessionsForPeriodResponse for returning reading sessions for period
type ReadingSessionsForPeriodResponse struct {
	Sessions []ReadingSessionResponse `json:"sessions"`
	Quantity int                      `json:"quantity"`
}

// ReadingSessionsUpdateRequest for updating reading session info
type ReadingSessionsUpdateRequest struct {
	BookID     int64         `json:"book_id" validate:"required,gt=0"`
	TelegramID int64         `json:"telegram_id" validate:"required,gt=0"`
	Duration   time.Duration `json:"duration" validate:"required,gt=0"`
	Pages      int64         `json:"pages" validate:"required,gt=0"`
	Date       time.Time     `json:"date" validate:"required"`
}

// BookListResponse for returning list of books
type BookListResponse struct {
	Books    []BookResponse `json:"books"`
	Quantity int            `json:"quantity"`
}

// UserStatsRequest for getting user stats
type UserStatsRequest struct {
	TelegramID int64 `json:"telegram_id" validate:"required,gt=0"`
}

// UserStatsResponse for returning user stats
type UserStatsResponse struct {
	TotalBooks     int64      `json:"total_books"`
	TotalPages     int64      `json:"total_pages"`
	BooksThisYear  int64      `json:"books_this_year"`
	PagesThisMonth int64      `json:"pages_this_month"`
	AvgPagesPerDay float64    `json:"avg_pages_per_day"`
	CurrentStreak  int64      `json:"current_streak"`
	LastReadDate   *time.Time `json:"last_read_date,omitempty"`
}
