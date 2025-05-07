package repositories

import (
	"pedika-layered-architecture/internal/models"
	"gorm.io/gorm"
)

type ViolenceCategoryRepository interface {
	FindAll() ([]models.ViolenceCategory, error)
	FindByID(id int64) (*models.ViolenceCategory, error)
	Create(category *models.ViolenceCategory) error
	Update(category *models.ViolenceCategory) error
	Delete(id int64) error
}

type violenceCategoryRepository struct {
	db *gorm.DB
}

func NewViolenceCategoryRepository(db *gorm.DB) ViolenceCategoryRepository {
	return &violenceCategoryRepository{db}
}

func (r *violenceCategoryRepository) FindAll() ([]models.ViolenceCategory, error) {
	var categories []models.ViolenceCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *violenceCategoryRepository) FindByID(id int64) (*models.ViolenceCategory, error) {
	var cat models.ViolenceCategory
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *violenceCategoryRepository) Create(cat *models.ViolenceCategory) error {
	return r.db.Create(cat).Error
}

func (r *violenceCategoryRepository) Update(cat *models.ViolenceCategory) error {
	return r.db.Save(cat).Error
}

func (r *violenceCategoryRepository) Delete(id int64) error {
	return r.db.Delete(&models.ViolenceCategory{}, id).Error
}
