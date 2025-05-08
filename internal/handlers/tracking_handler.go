package handlers

import (
	"mime/multipart"
	"strconv"

	"pedika-layered-architecture/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ReportTrackingHandler struct {
	service *services.ReportTrackingService
}

func NewReportTrackingHandler(s *services.ReportTrackingService) *ReportTrackingHandler {
	return &ReportTrackingHandler{service: s}
}

func (h *ReportTrackingHandler) Create(c *fiber.Ctx) error {
	noReg := c.FormValue("no_registrasi")
	keterangan := c.FormValue("keterangan")

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Gagal membaca form"})
	}

	files := form.File["documents"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Dokumen wajib diunggah"})
	}

	var openedFiles []multipart.File
	var filenames []string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuka file"})
		}
		defer src.Close()
		openedFiles = append(openedFiles, src)
		filenames = append(filenames, file.Filename)
	}

	if err := h.service.Create(noReg, keterangan, openedFiles, filenames); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Tracking berhasil dibuat"})
}

func (h *ReportTrackingHandler) GetByNoRegistrasi(c *fiber.Ctx) error {
	noReg := c.Params("noReg")

	trackings, err := h.service.GetByNoRegistrasi(noReg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(trackings) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Nomor Registrasi tidak ditemukan"})
	}

	return c.JSON(trackings)
}

func (h *ReportTrackingHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var body struct {
		Keterangan string `json:"keterangan"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}
	if err := h.service.Update(uint(id), body.Keterangan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Tracking berhasil diupdate"})
}

func (h *ReportTrackingHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Tracking berhasil dihapus"})
}
