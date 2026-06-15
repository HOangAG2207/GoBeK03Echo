package links

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/links"
	pkgutils "github.com/HOangAG2207/GoBeK03Echo/pkg/utils"
)

//go:generate mockery --name Service --filename links_service_mock.go --output ./mocks

type Service interface {
	ShortenLink(ctx context.Context, url string, codeLength int, exptime int64) (string, error)
	GetLink(ctx context.Context, code string) (string, error)
}
type service struct {
	repo       links.Repository
	randomCode pkgutils.CodeGenerator
}

func NewService(repo links.Repository, randomCode pkgutils.CodeGenerator) Service {
	return &service{
		repo:       repo,
		randomCode: randomCode,
	}
}
