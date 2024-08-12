package service

import (
	"context"
	"errors"
	"fmt"

	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/rogue0026/ssov2/internal/models"
	"github.com/rogue0026/ssov2/internal/ssoconfig"
	"github.com/rogue0026/ssov2/internal/storage"
	"github.com/rogue0026/ssov2/internal/token"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidLogin        = errors.New("login contains bad symbols")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrPasswordNotStrong   = errors.New("password length must be equal or greater than 8 symbols")
	ErrInvalidEmailAddress = errors.New("invalid email address")
	badChars               = "!@#$%^&*()+="
)

type UserModifier interface {
	CreateUser(ctx context.Context, u models.User) (int64, error)
	DeleteUser(ctx context.Context, login string, passHash []byte) error
}

type UserProvider interface {
	FetchUser(ctx context.Context, login string) (models.User, error)
}

type SSO struct {
	Config ssoconfig.SSOConfig
	UserModifier
	UserProvider
}

func New(modifier UserModifier, provider UserProvider) *SSO {
	service := SSO{
		UserModifier: modifier,
		UserProvider: provider,
	}

	return &service
}

func (s *SSO) RegisterNewUser(ctx context.Context, login string, password string, email string) (int64, error) {
	const fn = "internal.service.RegisterNewUser"
	if strings.ContainsAny(login, badChars) {
		return -1, ErrInvalidLogin
	}
	if utf8.RuneCountInString(password) < 8 {
		return -1, ErrPasswordNotStrong
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return -1, ErrInvalidEmailAddress
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}
	u := models.User{
		Login:        login,
		PasswordHash: hash,
		Email:        email,
	}
	userID, err := s.UserModifier.CreateUser(ctx, u)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}
	return userID, nil
}

func (s *SSO) LoginUser(ctx context.Context, login string, password string) (string, error) {
	const fn = "internal.service.LoginUser"
	usr, err := s.UserProvider.FetchUser(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", storage.ErrUserNotFound
		} else {
			return "", fmt.Errorf("%s: %w", fn, err)
		}
	}
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}
	tokenString, err := token.Generate(login, s.Config.TokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	return tokenString, nil
}
