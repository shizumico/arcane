package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	challengeHandlers "github.com/shizumico/arcane/cmd/api/internal/adapters/http/challenges"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/http/middleware"
	secretHandlers "github.com/shizumico/arcane/cmd/api/internal/adapters/http/secrets"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/repositories/redis"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/repositories/redis/challenges"
	"github.com/shizumico/arcane/cmd/api/internal/adapters/repositories/sqlite/secrets"
	challengesApplication "github.com/shizumico/arcane/cmd/api/internal/core/application/challenges"
	secretsApplication "github.com/shizumico/arcane/cmd/api/internal/core/application/secrets"
	"github.com/shizumico/arcane/pkg/logger"
	"github.com/shizumico/arcane/pkg/sqlite"
	"go.uber.org/zap"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	appLogger, err := logger.Init(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v\n", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting Arcane API server")

	db, err := sqlite.New(cfg.DbPath, cfg.MigrationsPath)
	if err != nil {
		appLogger.Fatal("Failed to connect to sqlite", zap.Error(err))
	}

	redisClient, err := redis.NewClient([]string{cfg.RedisHost + ":" + cfg.RedisPort}, cfg.RedisPassword)
	if err != nil {
		appLogger.Fatal("Failed to connect to redis", zap.Error(err))
	}

	challengeRepo := challenges.NewRepository(redisClient)
	secretQueryRepo := secrets.NewQueryRepository(db)
	secretCommandRepo := secrets.NewCommandRepository(db)

	challengeUseCase := challengesApplication.NewInteractor(challengeRepo)
	secretQueryUseCase := secretsApplication.NewQueryInteractor(secretQueryRepo)
	secretCommandUseCase := secretsApplication.NewCommandInteractor(secretCommandRepo)

	challengeHandlers := challengeHandlers.NewCommandHandlers(challengeUseCase)
	secretQueryHandlers := secretHandlers.NewQueryHandlers(secretQueryUseCase)
	secretCommandHandlers := secretHandlers.NewCommandHandlers(secretCommandUseCase, challengeUseCase)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "internal server error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				msg = e.Message
			}

			if code >= 500 {
				appLogger.Error("Server-side error",
					zap.Error(err),
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
				)
			} else {
				appLogger.Debug("Client-side error", zap.Int("status", code), zap.Error(err))
			}

			return c.Status(code).JSON(fiber.Map{"error": msg})
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(middleware.Authorization)

	api := app.Group("/api/v1")

	usernames := api.Group("/usernames")
	usernames.Get("/", secretQueryHandlers.ListUsernames)
	usernames.Get("/:username/services", secretQueryHandlers.ListServices)

	secrets := api.Group("/secrets")
	secrets.Post("/", secretCommandHandlers.Save)
	secrets.Get("/:username/:service", secretQueryHandlers.Get)

	api.Get("/challenge", challengeHandlers.Challenge)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	<-quit
	appLogger.Info("Shutting down server...")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited properly")
}
