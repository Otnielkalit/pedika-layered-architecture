package services

import (
	"errors"
	"io"
	"mime/multipart"
	"pedika-layered-architecture/internal/cloudinary"
	"pedika-layered-architecture/internal/models"
	"pedika-layered-architecture/internal/repositories"
)

type ReportTrackingService struct {
	repo       repositories.ReportTrackingRepository
	reportRepo repositories.ReportRepository
	cloudSvc   *cloudinary.CloudinaryService
}

func NewReportTrackingService(r repositories.ReportTrackingRepository, rr repositories.ReportRepository, c *cloudinary.CloudinaryService) *ReportTrackingService {
	return &ReportTrackingService{repo: r, reportRepo: rr, cloudSvc: c}
}

func (s *ReportTrackingService) Create(noReg, keterangan string, files []multipart.File, filenames []string) error {
	exists, err := s.reportRepo.ExistsByNoRegistrasi(noReg)
	if err != nil || !exists {
		return errors.New("no_registrasi tidak ditemukan")
	}

	readers := make([]io.Reader, len(files))
	for i, file := range files {
		readers[i] = file
	}
	urls, err := s.cloudSvc.UploadMultipleFiles(readers, filenames)
	if err != nil {
		return err
	}

	docMap := make(map[string]interface{})
	for i, url := range urls {
		docMap[filenames[i]] = url
	}
	tracking := models.ReportTracking{
		NoRegistrasi: noReg,
		Keterangan:   keterangan,
		Document:     docMap,
	}

	return s.repo.Create(tracking)
}

func (s *ReportTrackingService) GetByNoRegistrasi(noReg string) ([]models.ReportTracking, error) {
	return s.repo.GetByNoRegistrasi(noReg)
}

func (s *ReportTrackingService) Update(id uint, keterangan string) error {
	return s.repo.Update(id, keterangan)
}

func (s *ReportTrackingService) Delete(id uint) error {
	return s.repo.Delete(id)
}
