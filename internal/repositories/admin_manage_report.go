package repositories

import (
	"pedika-layered-architecture/internal/models"
	"time"

	"gorm.io/gorm"
)

type AdminManageReportRepository interface {
	GetAllReports() ([]models.Report, error)
	GetReportByRegistrationNumber(noRegistrasi string) (*models.Report, error)
	UpdateViewedReport(noRegistrasi string, userID uint, waktu *time.Time, status string) error
	UpdateProccessedReport(noRegistrasi string, userID uint, waktu *time.Time, status string) error
	UpdateCompletedReport(noRegistrasi string, userID uint, waktu *time.Time, status string) error
}

type adminManageReportRepository struct {
	db *gorm.DB
}

func NewAdminManageReportRepository(db *gorm.DB) AdminManageReportRepository {
	return &adminManageReportRepository{db: db}
}

func (r *adminManageReportRepository) GetAllReports() ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.
		Preload("User").
		Preload("ViolenceCategory").
		Order("created_at desc").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *adminManageReportRepository) GetReportByRegistrationNumber(noRegistrasi string) (*models.Report, error) {
	var report models.Report
	if err := r.db.Preload("ViolenceCategory").Preload("User").Where("no_registrasi = ?", noRegistrasi).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *adminManageReportRepository) UpdateViewedReport(noRegistrasi string, userID uint, waktu *time.Time, status string) error {
	return r.db.Model(&models.Report{}).
		Where("no_registrasi = ?", noRegistrasi).
		Updates(map[string]interface{}{
			"user_id_melihat": userID,
			"waktu_dilihat":   waktu,
			"status":          status,
		}).Error
}

// Update status menjadi Proses
func (r *adminManageReportRepository) UpdateProccessedReport(noRegistrasi string, userID uint, waktu *time.Time, status string) error {
	return r.db.Model(&models.Report{}).
		Where("no_registrasi = ?", noRegistrasi).
		Updates(map[string]interface{}{
			"user_id_melihat": userID,
			"waktu_diproses":  waktu,
			"status":          status,
		}).Error
}

// Update status menjadi Selesai
func (r *adminManageReportRepository) UpdateCompletedReport(noRegistrasi string, userID uint, waktu *time.Time, status string) error {
	return r.db.Model(&models.Report{}).
		Where("no_registrasi = ?", noRegistrasi).
		Updates(map[string]interface{}{
			"user_id_melihat": userID,
			"waktu_diproses":  waktu,
			"status":          status,
		}).Error
}
