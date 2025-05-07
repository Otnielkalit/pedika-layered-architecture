package services

import (
	"pedika-layered-architecture/internal/models"
	"pedika-layered-architecture/internal/repositories"
	"time"
)

type ReportService interface {
	GetReportsByUserID(userID uint) ([]models.Report, error)
	GetByNoRegistrasi(noReg string) (*models.Report, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) ReportService {
	return &reportService{reportRepo: repo}
}

func (s *reportService) GetReportsByUserID(userID uint) ([]models.Report, error) {
	return s.reportRepo.GetByUserID(userID)
}

func (s *reportService) GetByNoRegistrasi(noReg string) (*models.Report, error) {
	return s.reportRepo.GetByNoRegistrasi(noReg)
}

// CreateReport meng-handle logika untuk membuat laporan baru
func (s *reportService) CreateReport(userID uint, kategoriKekerasanID uint64, tanggalKejadian time.Time, buktiFile *fiber.File) (*models.Report, error) {
	// Create a new report entity
	report := models.Report{
		UserID:              userID,
		KategoriKekerasanID: kategoriKekerasanID,
		TanggalKejadian:     tanggalKejadian,
		Bukti:               *buktiFile,
	}

	
	return s.Repo.CreateReport(&report)
}
