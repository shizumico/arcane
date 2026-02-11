package challenges

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/challenges/dto"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/middleware"
	"github.com/shizumico/arcane/cmd/api/internal/core/application/challenges"
)

type CommandHandlers struct {
	useCase challenges.UseCase
}

func NewCommandHandlers(useCase challenges.UseCase) *CommandHandlers {
	return &CommandHandlers{useCase}
}

func (h *CommandHandlers) Challenge(c fiber.Ctx) error {
	pubkey := middleware.PubKeyFromCtx(c.Context())

	challenge, err := h.useCase.CreateChallenge(c.Context(), pubkey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Error: "failed to create challenge",
		})
	}

	return c.JSON(dto.Response{
		Challenge: challenge,
	})
}
