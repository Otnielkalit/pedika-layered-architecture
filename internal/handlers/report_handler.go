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


