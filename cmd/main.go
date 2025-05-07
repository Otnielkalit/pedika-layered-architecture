package main

import (
	"pedika-layered-architecture/config"
	"pedika-layered-architecture/database"
	"pedika-layered-architecture/internal/handlers"
	"pedika-layered-architecture/internal/repositories"
	"pedika-layered-architecture/internal/routes"
	"pedika-layered-architecture/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.InitDatabase()
	database.AutoMigrate()

	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(userService)

	routes.SetupRoutes(app, authHandler)

	app.Listen(":8087")
}
