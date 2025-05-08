package repositories

import (
	"pedika-layered-architecture/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetByUserID(userID uint) ([]models.Report, error)
	GetByNoRegistrasi(noReg string) (*models.Report, error)
	Create(report *models.Report) error
	GetLastNoRegistrasi() (string, error)
	Update(noReg string, report *models.Report) error
	Delete(noReg string) error
	Cancel(noReg string, alasan string) error
	ExistsByNoRegistrasi(noReg string) (bool, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetByUserID(userID uint) ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.
		Preload("User").
		Preload("ViolenceCategory").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *reportRepository) GetByNoRegistrasi(noReg string) (*models.Report, error) {
	var report models.Report
	if err := r.db.
		Preload("User").
		Preload("ViolenceCategory").
		Where("no_registrasi = ?", noReg).
		First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) GetLastNoRegistrasi() (string, error) {
	var lastReport models.Report
	err := r.db.Order("no_registrasi desc").First(&lastReport).Error
	if err == gorm.ErrRecordNotFound {
		return "000-DPMDPPA-IX-2000", nil
	}
	return lastReport.NoRegistrasi, err
}

func (r *reportRepository) Update(noReg string, report *models.Report) error {
	return r.db.Where("no_registrasi = ?", noReg).Save(report).Error
}

func (r *reportRepository) Delete(noReg string) error {
	return r.db.Where("no_registrasi = ?", noReg).Delete(&models.Report{}).Error
}

func (r *reportRepository) Cancel(noReg string, alasan string) error {
	now := time.Now()
	return r.db.Model(&models.Report{}).
		Where("no_registrasi = ?", noReg).
		Updates(map[string]interface{}{
			"alasan_dibatalkan": alasan,
			"waktu_dibatalkan":  now,
			"status":            "Dibatalkan",
		}).Error
}

func (r *reportRepository) ExistsByNoRegistrasi(noReg string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Report{}).Where("no_registrasi = ?", noReg).Count(&count).Error
	return count > 0, err
}