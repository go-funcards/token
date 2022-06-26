package token

import (
	"context"
	"fmt"
	"github.com/go-funcards/jwt"
	"github.com/google/uuid"
	"time"
)

var _ Service = (*service)(nil)

type Config struct {
	TokenType string        `yaml:"token_type" env:"TOKEN_TYPE" env-default:"Bearer"`
	TTL       time.Duration `yaml:"ttl" env:"TTL" env-default:"8760h"`
}

type Service interface {
	SessByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	SessByUser(ctx context.Context, user jwt.User) (Session, error)
}

type service struct {
	tokenType string
	ttl       time.Duration
	generator jwt.Generator
	storage   Storage
}

func New(cfg Config, generator jwt.Generator, storage Storage) *service {
	return &service{
		tokenType: cfg.TokenType,
		ttl:       cfg.TTL,
		generator: generator,
		storage:   storage,
	}
}

func (s *service) SessByRefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	user, err := s.storage.Get(ctx, refreshToken)
	if err != nil {
		return Session{}, fmt.Errorf("invalid refresh token: %w", err)
	}
	defer s.storage.Del(ctx, refreshToken)
	return s.SessByUser(ctx, user)
}

func (s *service) SessByUser(ctx context.Context, user jwt.User) (Session, error) {
	accessToken, err := s.generator.GenerateToken(user)
	if err != nil {
		return Session{}, fmt.Errorf("generate access token error: %w", err)
	}
	refreshToken := uuid.New().String()
	if err = s.storage.Set(ctx, refreshToken, user, s.ttl); err != nil {
		return Session{}, fmt.Errorf("save refresh token error: %w", err)
	}
	return Session{
		TokenType:    s.tokenType,
		ExpiresIn:    uint(s.ttl.Seconds()),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
