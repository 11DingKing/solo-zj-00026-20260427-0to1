package main

import (
	"log"
	"os"

	"rental-backend/config"
	"rental-backend/database"
	"rental-backend/middleware"
	"rental-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadConfig()

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(logger.New())

	app.Static("/uploads", "./uploads")

	routes.SetupPublicRoutes(app)
	app.Use(middleware.AuthMiddleware())
	routes.SetupProtectedRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
