package repository

import (
	"auth-service/internal/entities"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SessionPostgres struct {
	dbPool *pgxpool.Pool
}

func NewSessionPostgres(db *pgxpool.Pool) *SessionPostgres {
	return &SessionPostgres{dbPool: db}
}

func (db *SessionPostgres) AddSession(ctx context.Context, session *entities.Session) error {
	_, err := db.dbPool.Exec(ctx,
		`INSERT INTO sessions (session_id, login, create_date, expire_date, last_access_date)
		VALUES
		($1, $2, $3, $4, $5);`,
		session.SessionID, session.Login, session.CreateDate, session.ExpireDate, session.LastAccessDate)
	return err
}

func (db *SessionPostgres) DeleteSession(ctx context.Context, login string) error {
	_, err := db.dbPool.Exec(ctx,
		`DELETE FROM sessions WHERE login = $1;`, login)
	return err
}

func (db *SessionPostgres) GetLogin(ctx context.Context, sessionID string) (string, error) {
	var login string
	err := db.dbPool.QueryRow(
		ctx,
		"SELECT login FROM sessions WHERE session_id = $1;",
		sessionID).Scan(&login)
	return login, err
}

func (db *SessionPostgres) GetExpire(ctx context.Context, sessionID string) (time.Time, error) {
	var expireDate time.Time
	err := db.dbPool.QueryRow(
		ctx,
		"SELECT expire_date FROM sessions WHERE session_id = $1;",
		sessionID).Scan(&expireDate)
	return expireDate, err
}

func (db *SessionPostgres) IsExistsSession(ctx context.Context, login string) (bool, error) {
	var IsExistsSession bool
	err := db.dbPool.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM sessions WHERE login = $1);",
		login).Scan(&IsExistsSession)

	return IsExistsSession, err
}
