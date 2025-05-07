package services

import (
	"pedika-layered-architecture/internal/cloudinary"
	"pedika-layered-architecture/internal/models"
	"pedika-layered-architecture/internal/repositories"
	"fmt"
	"mime/multipart"
)

type ViolenceCategoryService struct {
	repo     repositories.ViolenceCategoryRepository
	cloudSvc *cloudinary.CloudinaryService
}

func NewViolenceCategoryService(r repositories.ViolenceCategoryRepository, c *cloudinary.CloudinaryService) *ViolenceCategoryService {
	return &ViolenceCategoryService{repo: r, cloudSvc: c}
}

func (s *ViolenceCategoryService) GetAll() ([]models.ViolenceCategory, error) {
	return s.repo.FindAll()
}

func (s *ViolenceCategoryService) GetByID(id int64) (*models.ViolenceCategory, error) {
	return s.repo.FindByID(id)
}

func (s *ViolenceCategoryService) Create(name string, file multipart.File, filename string) error {
	url, _, err := s.cloudSvc.UploadImage(file, filename)
	if err != nil {
		return err
	}
	cat := &models.ViolenceCategory{
		CategoryName: name,
		Image:        url,
	}
	return s.repo.Create(cat)
}

func (s *ViolenceCategoryService) Update(id int64, name string, file multipart.File, filename string) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	publicID := s.cloudSvc.GetPublicIDFromURL(cat.Image)
	if err := s.cloudSvc.DeleteImage(publicID); err != nil {
		fmt.Println("warning: gagal hapus image lama:", err)
	}
	newURL, _, err := s.cloudSvc.UploadImage(file, filename)
	if err != nil {
		return err
	}
	cat.CategoryName = name
	cat.Image = newURL
	return s.repo.Update(cat)
}

func (s *ViolenceCategoryService) Delete(id int64) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	publicID := s.cloudSvc.GetPublicIDFromURL(cat.Image)
	_ = s.cloudSvc.DeleteImage(publicID)

	return s.repo.Delete(id)
}
