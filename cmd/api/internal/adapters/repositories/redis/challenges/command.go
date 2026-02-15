package challenges

import (
	"context"
	"time"

	"github.com/redis/rueidis"
	"github.com/shizumico/arcane/cmd/api/internal/core/domain/challenges"
)

type Repository struct {
	redisClient rueidis.Client
}

func NewRepository(redisClient rueidis.Client) challenges.Repository {
	return &Repository{redisClient}
}

func (r *Repository) Save(ctx context.Context, pubkey, challenge string, ttl time.Duration) error {
	return r.redisClient.Do(ctx, r.redisClient.B().Set().
		Key("challenge:"+pubkey).
		Value(challenge).
		Ex(ttl).
		Build()).Error()
}

func (r *Repository) Get(ctx context.Context, pubkey string) (string, error) {
	challenge, err := r.redisClient.Do(ctx, r.redisClient.B().Get().
		Key("challenge:"+pubkey).
		Build()).ToString()
	return challenge, err
}

func (r *Repository) Delete(ctx context.Context, pubkey string) error {
	return r.redisClient.Do(ctx, r.redisClient.B().Del().
		Key("challenge:"+pubkey).
		Build()).Error()
}
