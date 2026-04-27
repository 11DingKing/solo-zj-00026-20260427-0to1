package handlers

import (
	"strconv"

	"rental-backend/database"
	"rental-backend/middleware"
	"rental-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SendMessageRequest struct {
	ReceiverID uuid.UUID `json:"receiver_id"`
	HouseID    uuid.UUID `json:"house_id"`
	Content    string    `json:"content"`
}

func GetMessages(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Message{}).Where("receiver_id = ?", user.ID)

	if unread := c.Query("unread"); unread == "true" {
		query = query.Where("is_read = ?", false)
	}

	var total int64
	query.Count(&total)

	var messages []models.Message
	query.Preload("Sender").Preload("House").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages)

	return c.JSON(fiber.Map{
		"data": messages,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetSentMessages(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Message{}).Where("sender_id = ?", user.ID)

	var total int64
	query.Count(&total)

	var messages []models.Message
	query.Preload("Receiver").Preload("House").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages)

	return c.JSON(fiber.Map{
		"data": messages,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func SendMessage(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)

	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ReceiverID == uuid.Nil || req.HouseID == uuid.Nil || req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Receiver ID, house ID, and content are required",
		})
	}

	if req.ReceiverID == user.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot send message to yourself",
		})
	}

	var house models.House
	if err := database.DB.First(&house, "id = ?", req.HouseID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "House not found",
		})
	}

	var receiver models.User
	if err := database.DB.First(&receiver, "id = ?", req.ReceiverID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Receiver not found",
		})
	}

	message := models.Message{
		SenderID:   user.ID,
		ReceiverID: req.ReceiverID,
		HouseID:    req.HouseID,
		Content:    req.Content,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	database.DB.Preload("Sender").Preload("Receiver").Preload("House").First(&message, message.ID)
	return c.Status(fiber.StatusCreated).JSON(message)
}

func ReplyMessage(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	parentIDStr := c.Params("id")
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	var parentMsg models.Message
	if err := database.DB.First(&parentMsg, "id = ?", parentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Message not found",
		})
	}

	if parentMsg.ReceiverID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only reply to messages sent to you",
		})
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content is required",
		})
	}

	message := models.Message{
		SenderID:   user.ID,
		ReceiverID: parentMsg.SenderID,
		HouseID:    parentMsg.HouseID,
		Content:    req.Content,
		ParentID:   &parentID,
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send reply",
		})
	}

	database.DB.Preload("Sender").Preload("Receiver").Preload("House").First(&message, message.ID)
	return c.Status(fiber.StatusCreated).JSON(message)
}

func MarkMessageAsRead(c *fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	messageIDStr := c.Params("id")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	var message models.Message
	if err := database.DB.First(&message, "id = ?", messageID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Message not found",
		})
	}

	if message.ReceiverID != user.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only mark messages sent to you as read",
		})
	}

	message.IsRead = true
	if err := database.DB.Save(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark message as read",
		})
	}

	return c.JSON(message)
}
