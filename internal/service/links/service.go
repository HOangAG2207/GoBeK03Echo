package links

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/links"
)

type Service interface {
	ShortenURL(ctx context.Context, url string, codeLength int, exptime int64) (string, error)
	GetURL(ctx context.Context, code string) (string, error)
}
type service struct {
	repo links.Repository
}

func NewService(repo links.Repository) Service {
	return &service{
		repo: repo,
	}
}
