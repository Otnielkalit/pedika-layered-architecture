package main

import (
	"pedika-layered-architecture/config"
	"pedika-layered-architecture/database"
	"pedika-layered-architecture/internal/cloudinary"
	"pedika-layered-architecture/internal/handlers"
	"pedika-layered-architecture/internal/repositories"
	"pedika-layered-architecture/internal/routes"
	"pedika-layered-architecture/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Init DB & migrate
	config.InitDatabase()
	database.AutoMigrate()
	db := config.DB

	// Init Cloudinary
	cloudSvc := cloudinary.NewCloudinaryService()

	// Auth
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(userService)

	// Violence Category
	categoryRepo := repositories.NewViolenceCategoryRepository(db)
	categoryService := services.NewViolenceCategoryService(categoryRepo, cloudSvc)
	categoryHandler := handlers.NewViolenceCategoryHandler(categoryService)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo, categoryRepo, cloudSvc)
	reportHandler := handlers.NewReportHandler(reportService)

	// Admin Manage Report
	adminReportRepo := repositories.NewAdminManageReportRepository(db)
	adminReportService := services.NewAdminManageReportService(adminReportRepo)
	adminReportHandler := handlers.NewAdminManageReportHandler(adminReportService)

	// Admin Manage Tracking Report

	// Admin Manage Report Tracking
	adminManageTrackingRepo := repositories.NewReportTrackingRepository(db)
	adminTrackingReportService := services.NewReportTrackingService(adminManageTrackingRepo, reportRepo, cloudSvc)
	adminTrackingReportHandler := handlers.NewReportTrackingHandler(adminTrackingReportService)

	// Routes
	routes.SetupRoutes(app, authHandler, categoryHandler, reportHandler, adminReportHandler, adminTrackingReportHandler)

	app.Listen(":8087")
}
