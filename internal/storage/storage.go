package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
)

type Storage interface {
	User() UserRepository
	//Dating() DatingRepository

	//Close() error
}
