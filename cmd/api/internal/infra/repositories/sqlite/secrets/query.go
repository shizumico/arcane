package secrets

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shizumico/arcane/cmd/api/internal/core/query/secrets"
)

type QueryRepository struct {
	db *sql.DB
}

func NewQueryRepository(db *sql.DB) secrets.Query {
	return &QueryRepository{db}
}

func (r *QueryRepository) ListUsernames(ctx context.Context, pubkey string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, ListUsernamesSQL, pubkey)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return r.getList(rows)
}

func (r *QueryRepository) ListServices(ctx context.Context, pubkey, username string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, ListServicesSQL, pubkey, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return r.getList(rows)
}

func (r *QueryRepository) Get(ctx context.Context, pubkey, username, service string) (secrets.SecretView, error) {
	var cipher, nonce string

	if err := r.db.QueryRowContext(ctx, GetSecretSQL, pubkey, username, service).
		Scan(&cipher, &nonce); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return secrets.SecretView{}, err
	}

	return toSecretView(username, service, cipher, nonce), nil
}

// -- Helpers --

func (r *QueryRepository) getList(rows *sql.Rows) ([]string, error) {
	entries := make([]string, 0)
	for rows.Next() {
		entry := ""
		if err := rows.Scan(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
