package routes

import (
	"pedika-layered-architecture/internal/handlers"
	"pedika-layered-architecture/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, categoryHandler *handlers.ViolenceCategoryHandler, reportHandler *handlers.ReportHandler) {
	api := app.Group("/api")

	// Auth
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

	// Admin
	admin := api.Group("/admin", middleware.AdminMiddleware)
	admin.Get("/violence-categories", categoryHandler.GetAll)
	admin.Get("/violence-categories/:id", categoryHandler.GetByID)
	admin.Post("/violence-categories", categoryHandler.Create)
	admin.Put("/violence-categories/:id", categoryHandler.Update)
	admin.Delete("/violence-categories/:id", categoryHandler.Delete)

	// Masyarakat
	masyarakat := api.Group("/masyarakat", middleware.MasyarakatMiddleware)
	masyarakat.Get("/reports", reportHandler.GetReportsByUserID)
	masyarakat.Get("/detail-report/:no_registrasi", reportHandler.GetReportByNoRegistrasi)
	masyarakat.Post("/report", reportHandler.CreateReport)
	masyarakat.Put("/report/:no_registrasi", reportHandler.UpdateReport)
}
