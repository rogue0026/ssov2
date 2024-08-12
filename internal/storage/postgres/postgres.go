package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rogue0026/ssov2/internal/models"
)

type Storage struct {
	*pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	const fn = "internal.storage.postgres.New"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	s := Storage{
		Pool: pool,
	}

	return &s, nil
}

func (s *Storage) Close() {
	s.Pool.Close()
}

func (s *Storage) CreateUser(ctx context.Context, u models.User) (int64, error) {
	const fn = "internal.storage.postgres.CreateUser"
	queryString := "INSERT INTO T_Users (login, password_hash, email) VALUES ($1, $2, $3);"
	_, err := s.Pool.Exec(ctx, queryString, u.Login, u.PasswordHash, u.Email)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	lastID, err := s.lastInsertUserID(ctx)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}
	return lastID, nil
}

func (s *Storage) DeleteUser(ctx context.Context, login string, passHash []byte) error {
	const fn = "internal.storage.postgres.DeleteUser"
	queryString := "DELETE FROM T_Users WHERE login = $1 AND password_hash = $2;"
	if _, err := s.Pool.Exec(ctx, queryString, login, passHash); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) FetchUser(ctx context.Context, login string) (models.User, error) {
	const fn = "internal.storage.postgres.FetchUser"
	queryString := "SELECT user_id, login, password_hash, email FROM T_Users WHERE login = $1;"
	row := s.Pool.QueryRow(ctx, queryString, login)
	var u models.User
	if err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.Email); err != nil {
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}
	return u, nil
}

func (s *Storage) lastInsertUserID(ctx context.Context) (int64, error) {
	const fn = "internal.storage.postgres.lastInsertUserID"
	queryString := "SELECT user_id FROM T_Users ORDER BY user_id DESC LIMIT 1"
	row := s.Pool.QueryRow(ctx, queryString)
	var lastID int64
	if err := row.Scan(&lastID); err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	return lastID, nil
}
