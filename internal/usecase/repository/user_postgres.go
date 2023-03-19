package repository

import (
	"auth-service/internal/entities"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserPostgres struct {
	dbPool *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{dbPool: db}
}

func (db *UserPostgres) AddUser(ctx context.Context, user *entities.User) error {
	_, err := db.dbPool.Exec(ctx,
		`INSERT INTO users (login, email, password_hash, phone_number)
		VALUES
		($1, $2, $3, $4)`,
		user.Login, user.Email, user.Password, user.PhoneNumber)

	return err
}

func (db *UserPostgres) IsValidLogin(ctx context.Context, login string) (bool, error) {
	var IsValidLogin bool
	err := db.dbPool.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1);",
		login).Scan(&IsValidLogin)

	return !IsValidLogin, err
}

func (db *UserPostgres) GetPassword(ctx context.Context, login string) (string, error) {
	var inputPassword string
	err := db.dbPool.QueryRow(
		ctx,
		"SELECT password_hash FROM users WHERE login = $1;",
		login).Scan(&inputPassword)
	return inputPassword, err
}

func (db *UserPostgres) GetUser(ctx context.Context, login string) (*entities.User, error) {
	var user entities.User
	user.Login = login
	err := db.dbPool.QueryRow(
		ctx,
		"SELECT login, email, phone_number FROM users WHERE login = $1;",
		login).Scan(&user.Password, &user.Email, &user.PhoneNumber)
	return &user, err
}
