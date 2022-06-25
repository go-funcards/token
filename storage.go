package token

import (
	"context"
	"github.com/go-funcards/jwt"
	"time"
)

type Storage interface {
	Set(ctx context.Context, refreshToken string, user jwt.User, expiration time.Duration) error
	Get(ctx context.Context, refreshToken string) (jwt.User, error)
	Del(ctx context.Context, refreshToken string)
}
