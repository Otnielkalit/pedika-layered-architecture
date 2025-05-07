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

