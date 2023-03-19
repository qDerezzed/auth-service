package repository

import (
	"auth-service/internal/entities"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type User interface {
	AddUser(ctx context.Context, user *entities.User) error
	IsValidLogin(ctx context.Context, login string) (bool, error)
	GetPassword(ctx context.Context, login string) (string, error)
	GetUser(ctx context.Context, login string) (*entities.User, error)
}

type Session interface {
	AddSession(ctx context.Context, session *entities.Session) error
	DeleteSession(ctx context.Context, login string) error
	GetLogin(ctx context.Context, sessionID string) (string, error)
	GetExpire(ctx context.Context, sessionID string) (time.Time, error)
	IsExistsSession(ctx context.Context, login string) (bool, error)
}

type Repository struct {
	User
	Session
}

func New(dbPool *pgxpool.Pool) *Repository {
	return &Repository{
		User:    NewUserPostgres(dbPool),
		Session: NewSessionPostgres(dbPool),
	}
}
