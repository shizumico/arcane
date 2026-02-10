package challenges

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"time"

	"github.com/shizumico/arcane/cmd/api/internal/core/domain/challenges"
	"github.com/shizumico/arcane/cmd/api/pkg/challenge"
	"github.com/shizumico/arcane/cmd/api/pkg/errors"
)

type UseCase interface {
	CreateChallenge(ctx context.Context, pubkey string) (string, error)
	VerifySignature(ctx context.Context, pubkey, signature string) error
}

type Interactor struct {
	repo challenges.Repository
}

func NewInteractor(repo challenges.Repository) UseCase {
	return &Interactor{repo}
}

func (c *Interactor) CreateChallenge(ctx context.Context, pubkey string) (string, error) {
	challenge := challenge.Generate()
	return challenge, c.repo.Save(ctx, pubkey, challenge, 5*time.Minute)
}

func (c *Interactor) VerifySignature(ctx context.Context, pubkey, signature string) error {
	pubBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return errors.ErrInvalidPubkeyFormat
	}

	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return errors.ErrInvalidSignatureFormat
	}

	challenge, err := c.repo.Get(ctx, pubkey)
	if err != nil {
		return err
	}

	if !ed25519.Verify(pubBytes, []byte(challenge), sigBytes) {
		return errors.ErrInvalidSignature
	}

	return c.repo.Delete(ctx, pubkey)
}
