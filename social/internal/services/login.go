package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
	"github.com/moeryomenko/highload-architect-otus/social/internal/repository"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidNickname = errors.New("invalid nickname")
)

// Login represents login service.
type Login struct {
	salt []byte

	repo *repository.Login
}

// NewLogin returns new instance of login service.
func NewLogin(cfg *config.Config, repo *repository.Login) *Login {
	return &Login{salt: []byte(cfg.Database.PasswordSalt), repo: repo}
}

// Check checks login password.
func (s *Login) Check(ctx context.Context, login *domain.Login) (*domain.Login, error) {
	l, err := s.repo.Get(ctx, login.Nickname)
	switch err {
	case nil:
		// it's ok.
	case repository.ErrNotFound:
		return nil, ErrInvalidNickname
	default:
		return nil, err
	}

	if s.encrypt(login.Password) != l.Password {
		return nil, ErrInvalidPassword
	}
	return l, nil
}

// EncryptSave encrypts login befor save to database.
func (s *Login) EncryptSave(ctx context.Context, login *domain.Login) error {
	login.Password = s.encrypt(login.Password)
	return s.repo.Save(ctx, login)
}

// encrypt encrypts user password.
func (s *Login) encrypt(password string) string {
	hash := hmac.New(sha256.New, s.salt)
	_, _ = hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
