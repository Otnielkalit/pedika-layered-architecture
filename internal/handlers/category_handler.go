package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pedika-layered-architecture/internal/services"
)

type ViolenceCategoryHandler struct {
	service *services.ViolenceCategoryService
}

func NewViolenceCategoryHandler(s *services.ViolenceCategoryService) *ViolenceCategoryHandler {
	return &ViolenceCategoryHandler{service: s}
}

func (h *ViolenceCategoryHandler) GetAll(c *fiber.Ctx) error {
	cats, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cats)
}

func (h *ViolenceCategoryHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	cat, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Kategori tidak ditemukan"})
	}
	return c.JSON(cat)
}

func (h *ViolenceCategoryHandler) Create(c *fiber.Ctx) error {
	name := c.FormValue("category_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama kategori wajib diisi",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gambar wajib diunggah",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuka file",
		})
	}
	defer src.Close()

	err = h.service.Create(name, src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Kategori berhasil dibuat",
	})
}

func (h *ViolenceCategoryHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	name := c.FormValue("category_name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama kategori wajib diisi",
		})
	}
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gambar wajib diunggah",
		})
	}
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuka file",
		})
	}
	defer src.Close()
	err = h.service.Update(id, name, src, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Kategori berhasil diperbarui",
	})
}

func (h *ViolenceCategoryHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Kategori berhasil dihapus"})
}
