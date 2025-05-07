package database

import (
	"log"
	"pedika-layered-architecture/config"
	"pedika-layered-architecture/internal/models"
)

func AutoMigrate() {
	err := config.DB.AutoMigrate(&models.User{}, &models.ViolenceCategory{}, &models.Report{})
	if err != nil {
		log.Fatalf("❌ Gagal migrasi database: %v", err)
	}
	log.Println("✅ Migrasi database berhasil")
}
