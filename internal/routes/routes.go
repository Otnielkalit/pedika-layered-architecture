package routes

import (
	"pedika-layered-architecture/config"
	"pedika-layered-architecture/internal/cloudinary"
	"pedika-layered-architecture/internal/handlers"
	"pedika-layered-architecture/internal/middleware"
	"pedika-layered-architecture/internal/repositories"
	"pedika-layered-architecture/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	api := app.Group("/api")

	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

	// 🔴 ADMIN ONLY
	admin := api.Group("/admin", middleware.AdminMiddleware)
	categoryHandler := handlers.NewViolenceCategoryHandler(
		services.NewViolenceCategoryService(
			repositories.NewViolenceCategoryRepository(config.DB),
			cloudinary.NewCloudinaryService(),
		),
	)
	admin.Get("/violence-categories", categoryHandler.GetAll)
	admin.Get("/violence-categories/:id", categoryHandler.GetByID)
	admin.Post("/violence-categories", categoryHandler.Create)
	admin.Put("/violence-categories/:id", categoryHandler.Update)
	admin.Delete("/violence-categories/:id", categoryHandler.Delete)

	// 🔵 MASYARAKAT ONLY
	reportHandler := handlers.NewReportHandler(
		services.NewReportService(
			repositories.NewReportRepository(config.DB),
		),
	)
	masyarakat := api.Group("/masyarakat", middleware.MasyarakatMiddleware)
	masyarakat.Get("/reports", reportHandler.GetReportsByUserID)
	masyarakat.Get("/detail-report/:no_registrasi", reportHandler.GetReportByNoRegistrasi)

}
