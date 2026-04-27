package handlers

import (
	"strconv"

	"rental-backend/database"
	"rental-backend/middleware"
	"rental-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetFavorites(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user.Role != models.RoleTenant {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only tenants can view favorites",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Favorite{}).Where("tenant_id = ?", user.ID)

	var total int64
	query.Count(&total)

	var favorites []models.Favorite
	query.Preload("House").Preload("House.Landlord").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&favorites)

	var houses []models.House
	for _, f := range favorites {
		if f.House != nil {
			houses = append(houses, *f.House)
		}
	}

	return c.JSON(fiber.Map{
		"data": houses,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func AddFavorite(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user.Role != models.RoleTenant {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only tenants can add favorites",
		})
	}

	houseIDStr := c.Params("houseId")
	houseID, err := uuid.Parse(houseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
	}

	var house models.House
	if err := database.DB.First(&house, "id = ? AND status = ?", houseID, models.StatusPublished).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "House not found",
		})
	}

	var existingFavorite models.Favorite
	if err := database.DB.Where("tenant_id = ? AND house_id = ?", user.ID, houseID).First(&existingFavorite).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Already favorited",
		})
	}

	favorite := models.Favorite{
		TenantID: user.ID,
		HouseID:  houseID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add favorite",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(favorite)
}

func RemoveFavorite(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user.Role != models.RoleTenant {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only tenants can remove favorites",
		})
	}

	houseIDStr := c.Params("houseId")
	houseID, err := uuid.Parse(houseIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
	}

	var favorite models.Favorite
	if err := database.DB.Where("tenant_id = ? AND house_id = ?", user.ID, houseID).First(&favorite).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Favorite not found",
		})
	}

	if err := database.DB.Delete(&favorite).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove favorite",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Favorite removed successfully",
	})
}
