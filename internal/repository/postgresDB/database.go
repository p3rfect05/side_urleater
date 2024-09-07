package postgresDB

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const linkExpireIn = 90 * 24 * time.Hour

type Storage struct {
	pgxPool      *pgxpool.Pool
	queryBuilder squirrel.StatementBuilderType
}

func NewStorage(pgxPool *pgxpool.Pool) *Storage {
	return &Storage{
		pgxPool:      pgxPool,
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (s *Storage) CreateUser(ctx context.Context, email string, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query, args, err := s.queryBuilder.Insert("users").
		Columns("email", "password_hash", "created_at").
		Values(email, passwordHash, time.Now().UTC().Format(time.RFC3339)).
		ToSql()

	if err != nil {
		return fmt.Errorf("CreateUser query error | %w", err)
	}

	_, err = s.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("CreateUser query error | %w", err)
	}

	return nil

}

func (s *Storage) ChangePassword(ctx context.Context, email string, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query, args, err := s.queryBuilder.
		Update("users").
		Set("password_hash", passwordHash).
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		return fmt.Errorf("ChangePassword query error | %w", err)
	}

	_, err = s.pgxPool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ChangePassword query error | %w", err)
	}

	return nil

}

func (s *Storage) GetUser(ctx context.Context, email string) (*User, error) {
	var user User

	query, args, err := s.queryBuilder.
		Select(
			"email",
			"password_hash",
			"urls_left",
		).
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("GetUser query error | %w", err)
	}

	err = s.pgxPool.QueryRow(ctx, query, args...).Scan(&user.Email, &user.PasswordHash, &user.UrlsLeft)
	if err != nil {
		return &User{}, fmt.Errorf("GetUser query error | %w", err)
	}
	return &user, nil
}

func (s *Storage) CreateShortLink(ctx context.Context, shortLink string, longLink string, userID int) (*Link, error) {
	var link Link
	expiresAt := time.Now().UTC().Add(linkExpireIn)
	query, args, err := s.queryBuilder.Insert("urls").
		Columns("short_url", "long_url", "created_at", "user_id", "expires_at").
		Values(shortLink, longLink, time.Now().UTC().Format(time.RFC3339), userID, expiresAt.Format(time.RFC3339)).
		Suffix("RETURNING short_url, long_url, expires_at").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("CreateShortLink query error | %w", err)
	}

	err = s.pgxPool.QueryRow(ctx, query, args...).Scan(
		&link.ShortUrl,
		&link.LongUrl,
		&link.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("CreateShortLink query error | %w", err)
	}

	return &link, nil
}

func (s *Storage) GetShortLink(ctx context.Context, shortLink string) (*Link, error) {
	var link Link

	query, args, err := s.queryBuilder.
		Select(
			"short_url",
			"long_url",
			"expires_at",
		).
		From("urls").
		Where(squirrel.Eq{"short_url": shortLink}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("GetShortLink query error | %w", err)
	}

	err = s.pgxPool.QueryRow(ctx, query, args...).Scan(&link.ShortUrl,
		&link.LongUrl,
		&link.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("GetShortLink query error | %w", err)
	}
	return &link, nil
}

func (s *Storage) DeleteShortLink(ctx context.Context, shortLink string) error {
	query, args, err := s.queryBuilder.
		Delete("urls").
		Where(squirrel.Eq{"short_url": shortLink}).
		ToSql()

	if err != nil {
		return fmt.Errorf("DeleteShortLink query error | %w", err)
	}

	_, err = s.pgxPool.Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("DeleteShortLink query error | %w", err)
	}

	return nil
}

func (s *Storage) ExtendShortLink(ctx context.Context, shortLink string, expiresAt time.Time) (*Link, error) {
	var link Link

	query, args, err := s.queryBuilder.
		Update("urls").
		Set("expires_at", expiresAt.Add(linkExpireIn).UTC().Format(time.RFC3339)).
		Where(squirrel.Eq{"short_url": shortLink}).
		Suffix("RETURNING short_url, long_url, expires_at").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("ExtendShortLink query error | %w", err)
	}

	err = s.pgxPool.QueryRow(ctx, query, args...).Scan(&link.ShortUrl,
		&link.LongUrl,
		&link.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("ExtendShortLink query error | %w", err)
	}

	return &link, nil
}
