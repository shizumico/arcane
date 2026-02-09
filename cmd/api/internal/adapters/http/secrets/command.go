package secrets

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/middleware"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/secrets/dto"
	"github.com/shizumico/arcane/cmd/api/internal/core/application/challenges"
	"github.com/shizumico/arcane/cmd/api/internal/core/application/secrets"
	errorsPkg "github.com/shizumico/arcane/cmd/api/pkg/errors"
)

type CommandHandlers struct {
	useCase          secrets.CommandUseCase
	challengeUseCase challenges.UseCase
}

func NewCommandHandlers(useCase secrets.CommandUseCase, challengeUseCase challenges.UseCase) *CommandHandlers {
	return &CommandHandlers{
		useCase:          useCase,
		challengeUseCase: challengeUseCase,
	}
}

func (h *CommandHandlers) Save(c fiber.Ctx) error {
	pubkey := middleware.PubKeyFromCtx(c.Context())

	signature := c.Get("Signature")
	if signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.SaveResponse{
			Error: "missing signature header",
		})
	}

	if err := h.challengeUseCase.VerifySignature(c.Context(), pubkey, signature); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.SaveResponse{
			Error: "failed to verify signature",
		})
	}

	req := dto.SaveSecretRequest{}
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.SaveResponse{
			Error: "invalid json",
		})
	}

	if err := h.useCase.Save(c.Context(), req.ToCommand(pubkey, signature)); err != nil {
		code, msg := h.mapErrorToStatusCode(err)
		return c.Status(code).JSON(dto.SaveResponse{
			Error: msg,
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *CommandHandlers) mapErrorToStatusCode(err error) (int, string) {
	switch {
	case errors.Is(err, errorsPkg.ErrInvalidPubkeyFormat):
		return fiber.StatusBadRequest, err.Error()
	case errors.Is(err, errorsPkg.ErrInvalidSignatureFormat):
		return fiber.StatusBadRequest, err.Error()
	case errors.Is(err, errorsPkg.ErrInvalidSignature):
		return fiber.StatusUnauthorized, err.Error()
	default:
		return fiber.StatusInternalServerError, "unknown internal error"
	}
}
