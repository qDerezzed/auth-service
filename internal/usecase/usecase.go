package usecase

import (
	"auth-service/internal/entities"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type AuthUseCase struct {
	repo AuthRepo
}

func New(r AuthRepo) *AuthUseCase {
	return &AuthUseCase{
		repo: r,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, user *entities.User) error {
	IsValidLogin, err := uc.repo.IsValidLogin(ctx, user.Login)
	if err != nil {
		return fmt.Errorf("AuthUseCase - Register - uc.IsValidLogin: %w", err)
	}
	if !IsValidLogin {
		return entities.ErrNotValidLogin
	}

	if err := uc.addUser(ctx,
		&entities.User{Login: user.Login,
			Email:       user.Email,
			Password:    user.Password,
			PhoneNumber: user.PhoneNumber,
		}); err != nil {
		return fmt.Errorf("AuthUseCase - Register - uc.AddUser: %w", err)
	}

	return nil
}

func (uc *AuthUseCase) CheckCreds(ctx context.Context, user *entities.User) (bool, error) {
	dbPassword, err := uc.repo.GetPassword(ctx, user.Login)
	if err != nil {
		return false, entities.ErrNotValidLoginOrPass
	}
	isValid, err := checkPass(dbPassword, user.Password)
	if err != nil {
		return false, fmt.Errorf("AuthUseCase - CheckCreds - uc.checkPass: %w", err)
	}
	return isValid, nil
}

func (uc *AuthUseCase) GenerateCookie(ctx context.Context, user *entities.User) (string, error) {
	// проверка на наличие в бд куки
	isExists, err := uc.repo.IsExistsSession(ctx, user.Login)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - GenerateCookie - uc.IsExistsSession: %w", err)
	}
	if isExists {
		if err := uc.repo.DeleteSession(ctx, user.Login); err != nil {
			return "", fmt.Errorf("AuthUseCase - GenerateCookie - uc.DeleteSession: %w", err)
		}
	}
	// создаем сессию
	sessionID, err := uc.addSession(ctx, user.Login)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - GenerateCookie - uc.AddSession: %w", err)
	}
	return sessionID, nil
}

func (uc *AuthUseCase) addUser(ctx context.Context, user *entities.User) error {
	salt := make([]byte, 8)
	rand.Read(salt)
	user.Password = base64.RawStdEncoding.EncodeToString(generatePasswordHash(salt, user.Password))
	err := uc.repo.AddUser(ctx, user)
	return err
}

func (uc *AuthUseCase) addSession(ctx context.Context, login string) (string, error) {
	sessionID := uc.generateSessionID()

	err := uc.repo.AddSession(ctx, &entities.Session{
		SessionID:      sessionID,
		Login:          login,
		CreateDate:     time.Now(),
		ExpireDate:     time.Now().Add(12 * time.Hour),
		LastAccessDate: time.Now(),
	})
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - addSession - uc.repo.AddSession: %w", err)
	}

	return sessionID, nil
}

func (uc *AuthUseCase) DeleteSession(ctx context.Context, login string) error {
	err := uc.repo.DeleteSession(ctx, login)
	if err != nil {
		return fmt.Errorf("AuthUseCase - DeleteSession - uc.repo.DeleteSession: %w", err)
	}
	return nil
}

func (uc *AuthUseCase) GetLogin(ctx context.Context, sessionID string) (string, error) {
	login, err := uc.repo.GetLogin(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - GetLogin - uc.repo.GetLogin: %w", err)
	}
	return login, nil
}

func (uc *AuthUseCase) GetUser(ctx context.Context, login string) (*entities.User, error) {
	user, err := uc.repo.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCase - GetUser - uc.repo.GetUser: %w", err)
	}
	return user, nil
}

func (uc *AuthUseCase) GetExpire(ctx context.Context, sessionID string) (time.Time, error) {
	expireDate, err := uc.repo.GetExpire(ctx, sessionID)
	if err != nil {
		return time.Time{}, fmt.Errorf("AuthUseCase - GetExpire - uc.repo.GetExpire: %w", err)
	}
	return expireDate, nil
}

func (uc *AuthUseCase) generateSessionID() string {
	sessionID := make([]byte, 32)
	rand.Read(sessionID)
	return base64.RawURLEncoding.EncodeToString(sessionID)
}

func generatePasswordHash(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func checkPass(passHash string, plainPassword string) (bool, error) {
	decodePassHash, err := base64.RawStdEncoding.Strict().DecodeString(passHash)
	if err != nil {
		return false, err
	}
	salt := decodePassHash[0:8]
	userPassHash := generatePasswordHash(salt, plainPassword)
	return bytes.Equal(userPassHash, decodePassHash), nil
}
