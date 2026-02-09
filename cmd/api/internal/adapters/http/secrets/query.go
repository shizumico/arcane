package secrets

import (
	"github.com/gofiber/fiber/v3"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/middleware"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/secrets/dto"
	"github.com/shizumico/arcane/cmd/api/internal/core/application/secrets"
)

type QueryHandlers struct {
	useCase secrets.QueryUseCase
}

func NewQueryHandlers(useCase secrets.QueryUseCase) *QueryHandlers {
	return &QueryHandlers{useCase}
}

func (h *QueryHandlers) ListUsernames(c fiber.Ctx) error {
	pubkey := middleware.PubKeyFromCtx(c.Context())

	usernames, err := h.useCase.ListUsernames(c.Context(), pubkey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ListUsernamesResponse{
			Error: "database error",
		})
	}

	if len(usernames) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.ListUsernamesResponse{
			Error: "usernames is empty",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.ListUsernamesResponse{
		Usernames: usernames,
	})
}

func (h *QueryHandlers) ListServices(c fiber.Ctx) error {
	pubkey := middleware.PubKeyFromCtx(c.Context())
	username := c.Params("username")

	services, err := h.useCase.ListServices(c.Context(), pubkey, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ListServicesResponse{
			Error: "database error",
		})
	}

	if len(services) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.ListServicesResponse{
			Error: "services is empty",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.ListServicesResponse{
		Services: services,
	})
}

func (h *QueryHandlers) Get(c fiber.Ctx) error {
	pubkey := middleware.PubKeyFromCtx(c.Context())
	username := c.Params("username")
	service := c.Params("service")

	secret, err := h.useCase.Get(c.Context(), pubkey, username, service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}

	if secret.Cipher == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "secret not found",
		})
	}

	return c.JSON(secret)
}
