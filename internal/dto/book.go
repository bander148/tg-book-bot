package dto

import "time"

type BookRequest struct {
	Pages       int64
	Description string
	Author      string
	Title       string
	PagesRead   int64
	StartDate   *time.Time
	EndDate     *time.Time
}

type UserRequest struct {
	TelegramID int64
	Username   string
}

type ReadingSession struct {
	Pages  int64
	BookID int64
}
