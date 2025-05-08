package handlers

import (
	"pedika-layered-architecture/internal/services"
	"pedika-layered-architecture/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminManageReportHandler struct {
	service services.AdminManageReportService
}

func NewAdminManageReportHandler(service services.AdminManageReportService) *AdminManageReportHandler {
	return &AdminManageReportHandler{service: service}
}

func (h *AdminManageReportHandler) GetAllReports(c *fiber.Ctx) error {
	reports, err := h.service.GetAllReports()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data laporan",
		})
	}
	return c.JSON(reports)
}

func (h *AdminManageReportHandler) GetReportByRegistrationNumber(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")

	report, err := h.service.GetReportByRegistrationNumber(noRegistrasi)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Laporan tidak ditemukan",
		})
	}

	return c.JSON(report)
}

func (h *AdminManageReportHandler) ViewReport(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")
	adminID, _ := utils.GetUserID(c)
	err := h.service.ViewReport(noRegistrasi, adminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memperbarui status laporan",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Laporan berhasil diperbarui menjadi sudah dilihat",
	})
}

func (h *AdminManageReportHandler) ProccessReport(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")
	adminID, _ := utils.GetUserID(c)

	err := h.service.ProccessReport(noRegistrasi, adminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memperbarui status laporan menjadi Proses",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Laporan berhasil diperbarui menjadi Proses",
	})
}

func (h *AdminManageReportHandler) CompleteReport(c *fiber.Ctx) error {
	noRegistrasi := c.Params("no_registrasi")
	adminID, _ := utils.GetUserID(c)
	err := h.service.CompleteReport(noRegistrasi, adminID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memperbarui status laporan menjadi Selesai",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Laporan berhasil diperbarui menjadi Selesai",
	})
}
