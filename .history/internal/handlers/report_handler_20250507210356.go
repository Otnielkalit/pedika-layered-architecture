package handlers

import (
	"pedika-layered-architecture/internal/services"
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
	userIDAny := c.Locals("user_id")
	userID, ok := userIDAny.(uint)
	if !ok {
		idStr := c.Locals("user_id").(string)
		parsedID, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}
		userID = uint(parsedID)
	}

	reports, err := h.service.GetReportsByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reports)
}

func (h *ReportHandler) GetReportByNoRegistrasi(c *fiber.Ctx) error {
	noReg := c.Params("no_registrasi")
	if noReg == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nomor registrasi wajib diisi"})
	}

	report, err := h.service.GetByNoRegistrasi(noReg)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Laporan tidak ditemukan"})
	}

	return c.JSON(report)
}

func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	// Ambil user ID dari konteks
	userID := c.Locals("user_id").(uint)

	// Ambil kategori kekerasan ID dari form dan validasi
	kategoriKekerasanID, err := strconv.ParseUint(c.FormValue("kategori_kekerasan_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Kategori kekerasan tidak valid"))
	}

	// Ambil tanggal kejadian dari form
	tanggalKejadian, err := time.Parse("2006-01-02", c.FormValue("tanggal_kejadian"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Tanggal kejadian tidak valid"))
	}

	// Ambil form file bukti
	buktiFile, err := c.FormFile("bukti")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse("Bukti tidak valid"))
	}

	// Buat laporan baru
	report, err := h.Service.CreateReport(userID, kategoriKekerasanID, tanggalKejadian, buktiFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse("Gagal membuat laporan"))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse("Laporan berhasil dibuat", report))
}