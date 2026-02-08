package secrets

import "context"

type Repository interface {
	Save(ctx context.Context, secret *Secret) error
}
