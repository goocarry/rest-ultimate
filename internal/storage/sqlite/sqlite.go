package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goocarry/rest-ultimate/internal/storage"
	"github.com/mattn/go-sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) User() storage.UserRepository {
	return s
}

func New(storagePath string) (storage.Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS user(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tg_user_id TEXT NOT NULL UNIQUE,
			first_name TEXT,
			second_name TEXT,
			workplace TEXT,
			phone TEXT,
			email TEXT,
			is_validated bool DEFAULT FALSE,
			interests text[]           
		);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) RegisterUser(user storage.User) (int64, error) {
	const op = "storage.sqlite.RegisterUser"

	stmt, err := s.db.Prepare("INSERT INTO main.user(tg_user_id) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(user.TgUserId)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, fmt.Errorf("user exists"))
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserByTelegramID(telegramID int64) (*storage.User, error) {
	return nil, nil
}
