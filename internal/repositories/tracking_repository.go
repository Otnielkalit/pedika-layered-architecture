package repositories

import (
	"log"
	"pedika-layered-architecture/internal/models"

	"gorm.io/gorm"
)

type ReportTrackingRepository interface {
	Create(tracking models.ReportTracking) error
	GetByNoRegistrasi(noReg string) ([]models.ReportTracking, error)
	Update(id uint, keterangan string) error
	Delete(id uint) error
}

type reportTrackingRepo struct {
	db *gorm.DB
}

func NewReportTrackingRepository(db *gorm.DB) ReportTrackingRepository {
	return &reportTrackingRepo{db: db}
}

func (r *reportTrackingRepo) Create(t models.ReportTracking) error {
	return r.db.Create(&t).Error
}

func (r *reportTrackingRepo) GetByNoRegistrasi(noReg string) ([]models.ReportTracking, error) {
	var trackings []models.ReportTracking
	err := r.db.Where("no_registrasi = ?", noReg).Find(&trackings).Error
	if err != nil {
		log.Println("Error fetching ReportTracking by no_registrasi:", err)
		return nil, err
	}
	return trackings, nil
}

func (r *reportTrackingRepo) Update(id uint, keterangan string) error {
	return r.db.Model(&models.ReportTracking{}).Where("id = ?", id).Update("keterangan", keterangan).Error
}

func (r *reportTrackingRepo) Delete(id uint) error {
	return r.db.Delete(&models.ReportTracking{}, id).Error
}
