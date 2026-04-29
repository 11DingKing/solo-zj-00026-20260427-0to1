package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"rental-backend/database"
	"rental-backend/middleware"
	"rental-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateHouseRequest struct {
	Title         string               `json:"title"`
	CommunityName string               `json:"community_name"`
	Address       string               `json:"address"`
	Longitude     float64              `json:"longitude"`
	Latitude      float64              `json:"latitude"`
	Price         int                  `json:"price"`
	Area          float64              `json:"area"`
	Bedrooms      int                  `json:"bedrooms"`
	LivingRooms   int                  `json:"living_rooms"`
	Bathrooms     int                  `json:"bathrooms"`
	Floor         int                  `json:"floor"`
	TotalFloors   int                  `json:"total_floors"`
	Orientation   models.Orientation   `json:"orientation"`
	Decoration    models.DecorationLevel `json:"decoration"`
	Facilities    []string             `json:"facilities"`
	MoveInDate    *time.Time           `json:"move_in_date"`
	MinLeaseTerm  int                  `json:"min_lease_term"`
	Description   string               `json:"description"`
	Images        []string             `json:"images"`
	Status        models.HouseStatus   `json:"status"`
}

func CreateHouse(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user.Role != models.RoleLandlord {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only landlords can create houses",
		})
	}

	var req CreateHouseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	house := models.House{
		LandlordID:    user.ID,
		Title:         req.Title,
		CommunityName: req.CommunityName,
		Address:       req.Address,
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
		Price:         req.Price,
		Area:          req.Area,
		Bedrooms:      req.Bedrooms,
		LivingRooms:   req.LivingRooms,
		Bathrooms:     req.Bathrooms,
		Floor:         req.Floor,
		TotalFloors:   req.TotalFloors,
		Orientation:   req.Orientation,
		Decoration:    req.Decoration,
		Facilities:    req.Facilities,
		MoveInDate:    req.MoveInDate,
		MinLeaseTerm:  req.MinLeaseTerm,
		Description:   req.Description,
		Images:        req.Images,
		Status:        req.Status,
	}

	if house.Status == "" {
		house.Status = models.StatusPublished
	}

	if err := database.DB.Create(&house).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create house",
		})
	}

	database.DB.Preload("Landlord").First(&house, house.ID)
	return c.Status(fiber.StatusCreated).JSON(house)
}

func GetHouses(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.House{}).Where("status = ?", models.StatusPublished)

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title ILIKE ? OR community_name ILIKE ? OR address ILIKE ? OR description ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if p, err := strconv.Atoi(minPrice); err == nil {
			query = query.Where("price >= ?", p)
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if p, err := strconv.Atoi(maxPrice); err == nil {
			query = query.Where("price <= ?", p)
		}
	}

	if minArea := c.Query("min_area"); minArea != "" {
		if a, err := strconv.ParseFloat(minArea, 64); err == nil {
			query = query.Where("area >= ?", a)
		}
	}

	if maxArea := c.Query("max_area"); maxArea != "" {
		if a, err := strconv.ParseFloat(maxArea, 64); err == nil {
			query = query.Where("area <= ?", a)
		}
	}

	if bedrooms := c.Query("bedrooms"); bedrooms != "" {
		if b, err := strconv.Atoi(bedrooms); err == nil {
			query = query.Where("bedrooms = ?", b)
		}
	}

	if orientation := c.Query("orientation"); orientation != "" {
		query = query.Where("orientation = ?", orientation)
	}

	if decoration := c.Query("decoration"); decoration != "" {
		query = query.Where("decoration = ?", decoration)
	}

	if facilities := c.Query("facilities"); facilities != "" {
		facilityList := strings.Split(facilities, ",")
		for _, f := range facilityList {
			query = query.Where("facilities @> ?", "["+f+"]")
		}
	}

	var total int64
	query.Count(&total)

	var houses []models.House
	query.Preload("Landlord").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&houses)

	return c.JSON(fiber.Map{
		"data": houses,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetHouse(c *fiber.Ctx) error {
	id := c.Params("id")
	houseID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
	}

	var house models.House
	if err := database.DB.Preload("Landlord").First(&house, "id = ?", houseID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "House not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get house",
		})
	}

	database.DB.Model(&house).Update("view_count", gorm.Expr("view_count + 1"))

	return c.JSON(house)
}

func UpdateHouse(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	id := c.Params("id")
	houseID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
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

	if house.LandlordID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not the owner of this house",
		})
	}

	var req CreateHouseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.CommunityName != "" {
		updates["community_name"] = req.CommunityName
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Longitude != 0 {
		updates["longitude"] = req.Longitude
	}
	if req.Latitude != 0 {
		updates["latitude"] = req.Latitude
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Area > 0 {
		updates["area"] = req.Area
	}
	if req.Bedrooms >= 0 {
		updates["bedrooms"] = req.Bedrooms
	}
	if req.LivingRooms >= 0 {
		updates["living_rooms"] = req.LivingRooms
	}
	if req.Bathrooms >= 0 {
		updates["bathrooms"] = req.Bathrooms
	}
	if req.Floor > 0 {
		updates["floor"] = req.Floor
	}
	if req.TotalFloors > 0 {
		updates["total_floors"] = req.TotalFloors
	}
	if req.Orientation != "" {
		updates["orientation"] = req.Orientation
	}
	if req.Decoration != "" {
		updates["decoration"] = req.Decoration
	}
	if req.Facilities != nil {
		updates["facilities"] = req.Facilities
	}
	if req.MoveInDate != nil {
		updates["move_in_date"] = req.MoveInDate
	}
	if req.MinLeaseTerm > 0 {
		updates["min_lease_term"] = req.MinLeaseTerm
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&house).Updates(updates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update house",
			})
		}
	}

	database.DB.Preload("Landlord").First(&house, house.ID)
	return c.JSON(house)
}

func DeleteHouse(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	id := c.Params("id")
	houseID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
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

	if house.LandlordID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not the owner of this house",
		})
	}

	if err := database.DB.Delete(&house).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete house",
		})
	}

	return c.JSON(fiber.Map{
		"message": "House deleted successfully",
	})
}

func UpdateHouseStatus(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	id := c.Params("id")
	houseID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
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

	if house.LandlordID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not the owner of this house",
		})
	}

	var req struct {
		Status models.HouseStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	validStatuses := map[models.HouseStatus]bool{
		models.StatusDraft:      true,
		models.StatusPublished:  true,
		models.StatusOffShelves: true,
		models.StatusRented:     true,
	}

	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	house.Status = req.Status
	if err := database.DB.Save(&house).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update house status",
		})
	}

	return c.JSON(house)
}

func UploadHouseImages(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	id := c.Params("id")
	houseID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid house ID",
		})
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

	if house.LandlordID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not the owner of this house",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid multipart form",
		})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No images provided",
		})
	}

	if len(house.Images)+len(files) > 9 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Maximum 9 images allowed",
		})
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create upload directory",
		})
	}

	var uploadedUrls []string

	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		filePath := filepath.Join(uploadDir, newFilename)

		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to open uploaded file",
			})
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save uploaded file",
			})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save uploaded file",
			})
		}

		uploadedUrls = append(uploadedUrls, "/uploads/"+newFilename)
	}

	house.Images = append(house.Images, uploadedUrls...)
	if err := database.DB.Save(&house).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update house images",
		})
	}

	return c.JSON(fiber.Map{
		"urls":   uploadedUrls,
		"images": house.Images,
	})
}

func GetMyHouses(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user.Role != models.RoleLandlord {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only landlords can view their houses",
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

	query := database.DB.Model(&models.House{}).Where("landlord_id = ?", user.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var houses []models.House
	query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&houses)

	return c.JSON(fiber.Map{
		"data": houses,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
