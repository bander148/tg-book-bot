package dto

// TODO : add json tags for future use in api
type BookCreateRequest struct {
	Pages       int64  `validate:"required,gt=0"`
	Description string `validate:"omitempty,max=2000"`
	Author      string `validate:"required,max=100"`
	Title       string `validate:"required,max=100"`
}
type AddPagesToBookProgressRequest struct {
	BookID    int64 `validate:"required,gt=0"`
	PagesRead int64 `validate:"required,gt=0,max=1000"`
}

type UserDTO struct {
	TelegramID int64  `validate:"required,gt=0"`
	Username   string `validate:"required,max=100"`
	//  FirstName string
	//  SecondName string
}

type ReadingSessionRequest struct {
	Pages  int64
	BookID int64
}
