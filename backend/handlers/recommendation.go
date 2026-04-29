package handlers

import (
	"fmt"
	"math"
	"strconv"

	"rental-backend/database"
	"rental-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSimilarHouses(c *fiber.Ctx) error {
	id := c.Params("id")
	houseID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	if limit < 1 || limit > 20 {
		limit = 8
	}

	var house models.House
	if err := database.DB.First(&house, "id = ?", houseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "House not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get house",
		})
	}

	priceTolerance := 0.2
	minPrice := int(float64(house.Price) * (1 - priceTolerance))
	maxPrice := int(float64(house.Price) * (1 + priceTolerance))

	var similarHouses []models.House

	query := database.DB.Model(&models.House{}).
		Where("id != ?", houseID).
		Where("status = ?", models.StatusPublished).
		Where(`(
			(community_name = ?) OR
			(bedrooms = ? AND living_rooms = ? AND bathrooms = ?)
		) AND price >= ? AND price <= ?`,
			house.CommunityName,
			house.Bedrooms, house.LivingRooms, house.Bathrooms,
			minPrice, maxPrice)

	query.Order(fmt.Sprintf(`
		CASE
			WHEN community_name = '%s' THEN 0
			ELSE 1
		END,
		ABS(price - %d) ASC
	`, house.CommunityName, house.Price)).
		Limit(limit).
		Preload("Landlord").
		Find(&similarHouses)

	scoredHouses := make([]scoredHouse, 0, len(similarHouses))
	for _, h := range similarHouses {
		score := calculateSimilarityScore(&house, &h)
		scoredHouses = append(scoredHouses, scoredHouse{
			House: h,
			Score: score,
		})
	}

	for i := range scoredHouses {
		for j := i + 1; j < len(scoredHouses); j++ {
			if scoredHouses[i].Score < scoredHouses[j].Score {
				scoredHouses[i], scoredHouses[j] = scoredHouses[j], scoredHouses[i]
			}
		}
	}

	result := make([]models.House, 0, len(scoredHouses))
	for _, sh := range scoredHouses {
		result = append(result, sh.House)
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

type scoredHouse struct {
	House models.House
	Score float64
}

func calculateSimilarityScore(source, target *models.House) float64 {
	score := 0.0

	if source.CommunityName == target.CommunityName {
		score += 40.0
	}

	if source.Bedrooms == target.Bedrooms {
		score += 20.0
	}
	if source.LivingRooms == target.LivingRooms {
		score += 10.0
	}
	if source.Bathrooms == target.Bathrooms {
		score += 10.0
	}

	priceDiff := math.Abs(float64(source.Price-target.Price)) / float64(source.Price)
	if priceDiff <= 0.05 {
		score += 20.0
	} else if priceDiff <= 0.1 {
		score += 15.0
	} else if priceDiff <= 0.15 {
		score += 10.0
	} else if priceDiff <= 0.2 {
		score += 5.0
	}

	areaDiff := math.Abs(source.Area - target.Area)
	if areaDiff <= 5 {
		score += 5.0
	} else if areaDiff <= 10 {
		score += 3.0
	}

	if source.Decoration == target.Decoration {
		score += 5.0
	}

	commonFacilities := 0
	for _, sf := range source.Facilities {
		for _, tf := range target.Facilities {
			if sf == tf {
				commonFacilities++
			}
		}
	}
	if len(source.Facilities) > 0 {
		score += float64(commonFacilities) * 2.0
	}

	return score
}
