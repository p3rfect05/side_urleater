package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"urleater/internal/repository/postgresDB"
)

type Storage interface {
	CreateUser(ctx context.Context, email string, password string) error
	ChangePassword(ctx context.Context, email string, password string) error
	GetUser(ctx context.Context, email string) (*postgresDB.User, error)
	CreateShortLink(ctx context.Context, shortLink string, longLink string, userID int) (*postgresDB.Link, error)
	GetShortLink(ctx context.Context, shortLink string) (*postgresDB.Link, error)
	DeleteShortLink(ctx context.Context, shortLink string) error
	ExtendShortLink(ctx context.Context, shortLink string, expiresAt time.Time) (*postgresDB.Link, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) LoginUser(ctx context.Context, email string, password string) error {
	user, err := s.storage.GetUser(ctx, email)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func (s *Service) RegisterUser(ctx context.Context, email string, password string) error {
	if !validatePassword(password) {
		return fmt.Errorf("password is too short (at least 8 symbols), got %d", len(password))
	}

	err := s.storage.CreateUser(ctx, email, password)
	if err != nil {
		return err
	}

	return nil
}

func validatePassword(password string) bool {
	return len(password) > 7 // TODO add more validation to password
}
