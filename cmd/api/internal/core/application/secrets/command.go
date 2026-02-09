package secrets

import (
	"context"

	"github.com/shizumico/arcane/cmd/api/internal/core/domain/secrets"
)

type SaveCommand struct {
	Pubkey    string
	Username  string
	Service   string
	Cipher    string
	Nonce     string
	Signature string
}

type CommandUseCase interface {
	Save(ctx context.Context, cmd SaveCommand) error
}

type CommandInteractor struct {
	repo secrets.Repository
}

func NewCommandInteractor(repo secrets.Repository) CommandUseCase {
	return &CommandInteractor{repo}
}

func (s *CommandInteractor) Save(ctx context.Context, cmd SaveCommand) error {
	return s.repo.Save(ctx, secrets.NewSecret(
		cmd.Pubkey,
		cmd.Username,
		cmd.Service,
		cmd.Cipher,
		cmd.Nonce,
	))
}
