package repositories

import (
	"pedika-layered-architecture/internal/models"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetByUserID(userID uint) ([]models.Report, error)
	GetByNoRegistrasi(noReg string) (*models.Report, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetByUserID(userID uint) ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepository) GetByNoRegistrasi(noReg string) (*models.Report, error) {
	var report models.Report
	if err := r.db.Where("no_registrasi = ?", noReg).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}


func (r *reportRepository) CreateReport(report *models.Report) (*models.Report, error) {
	if err := r.db.Create(report).Error; err != nil {
		return nil, err
	}
	return report, nil
}
