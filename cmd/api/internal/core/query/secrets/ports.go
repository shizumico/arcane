package secrets

import "context"

type Query interface {
	ListUsernames(ctx context.Context, pubkey string) ([]string, error)
	ListServices(ctx context.Context, pubkey, username string) ([]string, error)
	Get(ctx context.Context, pubkey, username, service string) (SecretView, error)
}
