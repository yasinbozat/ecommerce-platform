package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/config"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/handler"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/middleware"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/service"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/repository/postgres"
)

func main() {
	cfg := config.Load()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		panic(err)
	}
	userRepo := postgres.NewUserRepository(db)
	addressRepo := postgres.NewAddressRepository(db)

	userService := service.NewUserService(userRepo, addressRepo)
	authService := service.NewAuthService(userRepo, cfg)

	userHandler := handler.NewUserHandler(userService)
	addressHandler := handler.NewAddressHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()

	// internal route (no middleware)
	app.Get("/internal/auth/validate", authHandler.Validate)

	// api routes (middleware)
	api := app.Group("/api/v1", middleware.KeycloakMiddleware())

	// user routes
	api.Get("/users/me", userHandler.GetProfile)
	api.Put("/users/me", userHandler.UpdateProfile)

	// address routes
	api.Get("/users/me/addresses", addressHandler.List)
	api.Post("/users/me/adresses", addressHandler.Create)
	api.Put("/users/me/addresses/:id", addressHandler.Update)
	api.Delete("/users/me/addresses/:id", addressHandler.Delete)
	api.Patch("/users/me/addresses/:id/default", addressHandler.SetDefault)

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	prometheus := fiberprometheus.New("user-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Listen(":" + cfg.App.Port)
}
