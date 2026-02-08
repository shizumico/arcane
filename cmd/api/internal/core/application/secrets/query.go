package secrets

import (
	"context"

	"github.com/shizumico/arcane/cmd/api/internal/core/query/secrets"
)

type QueryUseCase interface {
	ListUsernames(ctx context.Context, pubkey string) ([]string, error)
	ListServices(ctx context.Context, pubkey, username string) ([]string, error)
	Get(ctx context.Context, pubkey, username, service string) (secrets.SecretView, error)
}

type QueryInteractor struct {
	query secrets.Query
}

func NewQueryInteractor(query secrets.Query) QueryUseCase {
	return &QueryInteractor{query}
}

func (s *QueryInteractor) ListUsernames(ctx context.Context, pubkey string) ([]string, error) {
	return s.query.ListUsernames(ctx, pubkey)
}

func (s *QueryInteractor) ListServices(ctx context.Context, pubkey, username string) ([]string, error) {
	return s.query.ListServices(ctx, pubkey, username)
}

func (s *QueryInteractor) Get(ctx context.Context, pubkey, username, service string) (secrets.SecretView, error) {
	return s.query.Get(ctx, pubkey, username, service)
}
