package services

import (
	"pedika-layered-architecture/internal/models"
	"pedika-layered-architecture/internal/repositories"
	"time"
)

type AdminManageReportService interface {
	GetAllReports() ([]models.Report, error)
	GetReportByRegistrationNumber(noRegistrasi string) (*models.Report, error)
	ViewReport(noRegistrasi string, userID uint) error
	ProccessReport(noRegistrasi string, userID uint) error
	CompleteReport(noRegistrasi string, userID uint) error
}

type adminManageReportService struct {
	repo repositories.AdminManageReportRepository
}

func NewAdminManageReportService(repo repositories.AdminManageReportRepository) AdminManageReportService {
	return &adminManageReportService{repo: repo}
}

func (s *adminManageReportService) GetAllReports() ([]models.Report, error) {
	return s.repo.GetAllReports()
}

func (s *adminManageReportService) GetReportByRegistrationNumber(noRegistrasi string) (*models.Report, error) {
	return s.repo.GetReportByRegistrationNumber(noRegistrasi)
}

func (s *adminManageReportService) ViewReport(noRegistrasi string, userID uint) error {
	now := time.Now()
	status := "Sudah dilihat"
	return s.repo.UpdateViewedReport(noRegistrasi, userID, &now, status)
}

// Proses laporan (status Proses)
func (s *adminManageReportService) ProccessReport(noRegistrasi string, userID uint) error {
	now := time.Now()
	status := "Proses"
	return s.repo.UpdateProccessedReport(noRegistrasi, userID, &now, status)
}

// Selesaikan laporan (status Selesai)
func (s *adminManageReportService) CompleteReport(noRegistrasi string, userID uint) error {
	now := time.Now()
	status := "Selesai"
	return s.repo.UpdateCompletedReport(noRegistrasi, userID, &now, status)
}
