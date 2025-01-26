package storage

type User struct {
	ID        int64
	TgUserId  string
	FirstName string
	LastName  string
	Workplace string
	Email     string
	Phone     string
	IsValid   bool
	Interests []string
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=UserRepository
type UserRepository interface {
	RegisterUser(user User) (int64, error)
	GetUserByTelegramID(telegramID int64) (*User, error)
}
