package services

import (
	"fmt"
	"mime/multipart"
	"pedika-layered-architecture/internal/cloudinary"
	"pedika-layered-architecture/internal/models"
	"pedika-layered-architecture/internal/repositories"
	"pedika-layered-architecture/internal/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type ReportService interface {
	GetReportsByUserID(userID uint) ([]models.Report, error)
	GetByNoRegistrasi(noReg string) (*models.Report, error)
	CreateReport(userID uint, req *CreateReportRequest) (*models.Report, error)
	UpdateReport(userID uint, req *UpdateReportRequest) error
	DeleteReport(noRegistrasi string) error
	CancelReport(noReg string, alasan string) error
}

type CreateReportRequest struct {
	KategoriKekerasanID uint                    `json:"kategori_kekerasan_id" validate:"required"`
	TanggalKejadian     string                  `json:"tanggal_kejadian" validate:"required"`
	AlamatTKP           string                  `json:"alamat_tkp" validate:"required"`
	AlamatDetailTKP     string                  `json:"alamat_detail_tkp" validate:"required"`
	KronologisKasus     string                  `json:"kronologis_kasus" validate:"required"`
	KategoriLokasiKasus string                  `json:"kategori_lokasi_kasus" validate:"required"`
	Dokumentasi         []*multipart.FileHeader `json:"-"`
}

type UpdateReportRequest struct {
	NoRegistrasi        string                  `json:"-"`
	KategoriKekerasanID uint                    `json:"kategori_kekerasan_id"`
	TanggalKejadian     string                  `json:"tanggal_kejadian"`
	AlamatTKP           string                  `json:"alamat_tkp"`
	AlamatDetailTKP     string                  `json:"alamat_detail_tkp"`
	KronologisKasus     string                  `json:"kronologis_kasus"`
	KategoriLokasiKasus string                  `json:"kategori_lokasi_kasus"`
	Dokumentasi         []*multipart.FileHeader `json:"-"`
}
type reportService struct {
	reportRepo        repositories.ReportRepository
	categoryRepo      repositories.ViolenceCategoryRepository
	cloudinaryService *cloudinary.CloudinaryService
}

func NewReportService(
	reportRepo repositories.ReportRepository,
	categoryRepo repositories.ViolenceCategoryRepository,
	cloudinaryService *cloudinary.CloudinaryService,
) ReportService {
	return &reportService{
		reportRepo:        reportRepo,
		categoryRepo:      categoryRepo,
		cloudinaryService: cloudinaryService,
	}
}

func (s *reportService) GetReportsByUserID(userID uint) ([]models.Report, error) {
	return s.reportRepo.GetByUserID(userID)
}

func (s *reportService) GetByNoRegistrasi(noReg string) (*models.Report, error) {
	return s.reportRepo.GetByNoRegistrasi(noReg)
}

func (s *reportService) CreateReport(userID uint, req *CreateReportRequest) (*models.Report, error) {
	category, err := s.categoryRepo.FindByID(int64(req.KategoriKekerasanID))
	if err != nil || category == nil {
		return nil, fmt.Errorf("kategori kekerasan tidak valid")
	}
	noReg, err := s.generateNoRegistrasi()
	if err != nil {
		return nil, fmt.Errorf("gagal generate nomor registrasi")
	}
	var dokumentasiURLs []string
	for _, file := range req.Dokumentasi {
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membuka file")
		}
		defer src.Close()

		url, _, err := s.cloudinaryService.UploadImage(src, fmt.Sprintf("report_%s_%s", noReg, file.Filename))
		if err != nil {
			return nil, fmt.Errorf("gagal upload dokumentasi")
		}
		dokumentasiURLs = append(dokumentasiURLs, url)
	}
	tanggalKejadian, err := time.Parse("2006-01-02", req.TanggalKejadian)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid (YYYY-MM-DD)")
	}
	report := &models.Report{
		NoRegistrasi:        noReg,
		UserID:              userID,
		KategoriKekerasanID: req.KategoriKekerasanID,
		TanggalKejadian:     tanggalKejadian,
		AlamatTKP:           req.AlamatTKP,
		AlamatDetailTKP:     req.AlamatDetailTKP,
		KronologisKasus:     req.KronologisKasus,
		KategoriLokasiKasus: req.KategoriLokasiKasus,
		Dokumentasi:         datatypes.JSONMap{"urls": dokumentasiURLs},
		Status:              "Laporan Masuk",
		TanggalPelaporan:    time.Now(),
	}
	if err := s.reportRepo.Create(report); err != nil {
		return nil, fmt.Errorf("gagal menyimpan laporan")
	}

	return report, nil
}

func (s *reportService) UpdateReport(userID uint, req *UpdateReportRequest) error {
	report, err := s.reportRepo.GetByNoRegistrasi(req.NoRegistrasi)
	if err != nil {
		return fmt.Errorf("laporan tidak ditemukan")
	}

	category, err := s.categoryRepo.FindByID(int64(req.KategoriKekerasanID))
	if err != nil || category == nil {
		return fmt.Errorf("kategori kekerasan tidak valid")
	}
	tanggalKejadian, err := time.Parse("2006-01-02", req.TanggalKejadian)
	if err != nil {
		return fmt.Errorf("format tanggal tidak valid (YYYY-MM-DD)")
	}
	var dokumentasiURLs []string
	if req.Dokumentasi != nil {
		for _, file := range req.Dokumentasi {
			src, err := file.Open()
			if err != nil {
				return fmt.Errorf("gagal membuka file")
			}
			defer src.Close()

			url, _, err := s.cloudinaryService.UploadImage(src, fmt.Sprintf("report_%s_%s", req.NoRegistrasi, file.Filename))
			if err != nil {
				return fmt.Errorf("gagal upload dokumentasi")
			}
			dokumentasiURLs = append(dokumentasiURLs, url)
		}
		report.Dokumentasi = datatypes.JSONMap{"urls": dokumentasiURLs}
	}
	report.KategoriKekerasanID = req.KategoriKekerasanID
	report.TanggalKejadian = tanggalKejadian
	report.AlamatTKP = req.AlamatTKP
	report.AlamatDetailTKP = req.AlamatDetailTKP
	report.KronologisKasus = req.KronologisKasus
	report.KategoriLokasiKasus = req.KategoriLokasiKasus
	if err := s.reportRepo.Update(req.NoRegistrasi, report); err != nil {
		return fmt.Errorf("gagal memperbarui laporan")
	}

	return nil
}

func (s *reportService) DeleteReport(noRegistrasi string) error {
	report, err := s.reportRepo.GetByNoRegistrasi(noRegistrasi)
	if err != nil {
		return fmt.Errorf("laporan tidak ditemukan")
	}

	return s.reportRepo.Delete(report.NoRegistrasi)
}

func (s *reportService) CancelReport(noReg string, alasan string) error {
	if strings.TrimSpace(alasan) == "" {
		return fmt.Errorf("alasan pembatalan wajib diisi")
	}
	report, err := s.reportRepo.GetByNoRegistrasi(noReg)
	if err != nil {
		return fmt.Errorf("laporan tidak ditemukan")
	}
	if report.Status == "Dibatalkan" {
		return fmt.Errorf("laporan sudah dibatalkan sebelumnya")
	}
	return s.reportRepo.Cancel(noReg, alasan)
}

func (s *reportService) generateNoRegistrasi() (string, error) {
	lastNoReg, err := s.reportRepo.GetLastNoRegistrasi()
	if err != nil {
		return "", err
	}
	parts := strings.Split(lastNoReg, "-")
	seq, _ := strconv.Atoi(parts[0])
	newSeq := fmt.Sprintf("%03d", seq+1)
	monthRoman := utils.ConvertMonthToRoman(time.Now().Month())

	return fmt.Sprintf("%s-DPMDPPA-%s-%d",
		newSeq,
		monthRoman,
		time.Now().Year(),
	), nil
}
