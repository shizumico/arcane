package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type ctxKey string

const pubkeyCtxKey ctxKey = "pubkey"

func Authorization(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization scheme")
	}

	pubkey := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if pubkey == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "empty bearer token")
	}

	ctx := context.WithValue(c.Context(), pubkeyCtxKey, pubkey)
	c.SetContext(ctx)

	return c.Next()
}

// -- Helpers --

func PubKeyFromCtx(ctx context.Context) string {
	pubkey, ok := ctx.Value(pubkeyCtxKey).(string)
	if !ok {
		panic("auth middleware not applied")
	}
	return pubkey
}
