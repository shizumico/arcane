package challenges

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, pubkey, challenge string, ttl time.Duration) error
	Get(ctx context.Context, pubkey string) (string, error)
	Delete(ctx context.Context, pubkey string) error
}
