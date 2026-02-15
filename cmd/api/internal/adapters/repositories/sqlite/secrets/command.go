package secrets

import (
	"context"
	"database/sql"

	"github.com/shizumico/arcane/cmd/api/internal/core/domain/secrets"
)

type CommandRepository struct {
	db *sql.DB
}

func NewCommandRepository(db *sql.DB) secrets.Repository {
	return &CommandRepository{db}
}

func (r *CommandRepository) Save(ctx context.Context, secret *secrets.Secret) error {
	_, err := r.db.ExecContext(ctx, UpsertSecretSQL, secret.PubKey, secret.Username, secret.Service, secret.Cipher, secret.Nonce)
	return err
}
