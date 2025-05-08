package handlers

import (
	"fmt"
	"pedika-layered-architecture/internal/services"
	"pedika-layered-architecture/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetReportsByUserID(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	reports, err := h.service.GetReportsByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(reports)
}

func (h *ReportHandler) GetReportByNoRegistrasi(c *fiber.Ctx) error {
	noReg := c.Params("no_registrasi")
	report, err := h.service.GetByNoRegistrasi(noReg)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Laporan tidak ditemukan"})
	}
	return c.JSON(report)
}

func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}
	requiredFields := map[string]string{
		"kategori_kekerasan_id": form.Value["kategori_kekerasan_id"][0],
		"tanggal_kejadian":      form.Value["tanggal_kejadian"][0],
		"alamat_tkp":            form.Value["alamat_tkp"][0],
		"alamat_detail_tkp":     form.Value["alamat_detail_tkp"][0],
		"kronologis_kasus":      form.Value["kronologis_kasus"][0],
		"kategori_lokasi_kasus": form.Value["kategori_lokasi_kasus"][0],
	}
	for key, val := range requiredFields {
		if val == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("%s wajib diisi", key)})
		}
	}
	kategoriID, _ := strconv.Atoi(requiredFields["kategori_kekerasan_id"])
	noRegistrasi, err := h.service.CreateReport(userID, &services.CreateReportRequest{
		KategoriKekerasanID: uint(kategoriID),
		TanggalKejadian:     requiredFields["tanggal_kejadian"],
		AlamatTKP:           requiredFields["alamat_tkp"],
		AlamatDetailTKP:     requiredFields["alamat_detail_tkp"],
		KronologisKasus:     requiredFields["kronologis_kasus"],
		KategoriLokasiKasus: requiredFields["kategori_lokasi_kasus"],
		Dokumentasi:         form.File["dokumentasi"],
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "Laporan berhasil dikirim",
		"no_registrasi": noRegistrasi,
	})
}

func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	noReg := c.Params("no_registrasi")
	if noReg == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no_registrasi kosong"})
	}
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "format form tidak valid"})
	}
	get := func(key string) string {
		if val, ok := form.Value[key]; ok && len(val) > 0 {
			return val[0]
		}
		return ""
	}
	kategoriID, err := strconv.ParseUint(get("kategori_kekerasan_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "kategori_kekerasan_id tidak valid"})
	}
	req := &services.UpdateReportRequest{
		NoRegistrasi:        noReg,
		KategoriKekerasanID: uint(kategoriID),
		TanggalKejadian:     get("tanggal_kejadian"),
		AlamatTKP:           get("alamat_tkp"),
		AlamatDetailTKP:     get("alamat_detail_tkp"),
		KronologisKasus:     get("kronologis_kasus"),
		KategoriLokasiKasus: get("kategori_lokasi_kasus"),
		Dokumentasi:         form.File["dokumentasi"],
	}

	err = h.service.UpdateReport(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Laporan berhasil diperbarui"})
}

func (h *ReportHandler) DeleteReport(c *fiber.Ctx) error {
	noReg := c.Params("no_registrasi")
	err := h.service.DeleteReport(noReg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Laporan berhasil dihapus",
	})
}

func (h *ReportHandler) CancelReport(c *fiber.Ctx) error {
	noReg := c.Params("no_registrasi")
	type CancelPayload struct {
		Alasan string `json:"alasan"`
	}
	var payload CancelPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "format request tidak valid"})
	}
	if payload.Alasan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "alasan wajib diisi"})
	}

	err := h.service.CancelReport(noReg, payload.Alasan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Laporan berhasil dibatalkan"})
}
