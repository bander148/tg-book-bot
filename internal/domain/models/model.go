package model

import "time"

// Book
type Book struct {
	ID          int64
	UserID      int64
	Pages       int64
	Description string
	Author      string
	Title       string
	StartDate   time.Time
	EndDate     *time.Time
	PagesRead   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User
type User struct {
	TelegramID int64
	Username   string
	CreatedAt  time.Time
}

// ReadingSession for inromation about reading at one day
type ReadingSession struct {
	ID        int64
	BookID    int64
	UserID    int64
	Date      time.Time
	Pages     int64
	CreatedAt time.Time
}

// UserStats maybe implement at the end project (Not necessarily)
type UserStats struct {
	TotalBooks     int64
	TotalPages     int64
	BooksThisYear  int64
	PagesThisMonth int64
	AvgPagesPerDay float64
	CurrentStreak  int64
	LastReadDate   *time.Time
}
