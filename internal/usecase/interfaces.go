package usecase

import (
	"auth-service/internal/entities"
	"context"
	"time"
)

type (
	Auth interface {
		Register(ctx context.Context, user *entities.User) error
		CheckCreds(ctx context.Context, user *entities.User) (bool, error)
		GenerateCookie(ctx context.Context, user *entities.User) (string, error)
		DeleteSession(ctx context.Context, login string) error
		GetLogin(ctx context.Context, sessionID string) (string, error)
		GetUser(ctx context.Context, login string) (*entities.User, error)
		GetExpire(ctx context.Context, sessionID string) (time.Time, error)
	}

	AuthRepo interface {
		AddUser(ctx context.Context, user *entities.User) error
		AddSession(ctx context.Context, session *entities.Session) error
		DeleteSession(ctx context.Context, login string) error
		IsValidLogin(ctx context.Context, login string) (bool, error)
		GetPassword(ctx context.Context, login string) (string, error)
		GetLogin(ctx context.Context, sessionID string) (string, error)
		GetUser(ctx context.Context, login string) (*entities.User, error)
		GetExpire(ctx context.Context, sessionID string) (time.Time, error)
		IsExistsSession(ctx context.Context, login string) (bool, error)
	}
)
