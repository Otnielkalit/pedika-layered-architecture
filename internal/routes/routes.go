package routes

import (
	"pedika-layered-architecture/internal/handlers"
	"pedika-layered-architecture/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler,
	categoryHandler *handlers.ViolenceCategoryHandler,
	reportHandler *handlers.ReportHandler,
	adminManageReportHandler *handlers.AdminManageReportHandler,
	adminManageReportTrackingHandler *handlers.ReportTrackingHandler) {
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

	admin.Get("/reports", adminManageReportHandler.GetAllReports)
	admin.Get("/reports/:no_registrasi", adminManageReportHandler.GetReportByRegistrationNumber)
	admin.Put("/reports/view/:no_registrasi", adminManageReportHandler.ViewReport)
	admin.Put("/reports/proccess/:no_registrasi", adminManageReportHandler.ProccessReport)
	admin.Put("/reports/complete/:no_registrasi", adminManageReportHandler.CompleteReport)

	// Admin Report Tracking Routes
	admin.Post("/report-tracking", adminManageReportTrackingHandler.Create)
	admin.Get("/report-tracking/:no_registrasi", adminManageReportTrackingHandler.GetByNoRegistrasi)
	admin.Put("/report-tracking/:id", adminManageReportTrackingHandler.Update)
	admin.Delete("/report-tracking/:id", adminManageReportTrackingHandler.Delete)

	// Masyarakat
	masyarakat := api.Group("/masyarakat", middleware.MasyarakatMiddleware)
	masyarakat.Get("/reports", reportHandler.GetReportsByUserID)
	masyarakat.Get("/detail-report/:no_registrasi", reportHandler.GetReportByNoRegistrasi)
	masyarakat.Post("/report", reportHandler.CreateReport)
	masyarakat.Put("/report/:no_registrasi", reportHandler.UpdateReport)
	masyarakat.Delete("/report/:no_registrasi", reportHandler.DeleteReport)
	masyarakat.Put("/report/:no_registrasi/cancel", reportHandler.CancelReport)

}
