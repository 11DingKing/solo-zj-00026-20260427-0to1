package routes

import (
	"rental-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	houses := api.Group("/houses")
	houses.Get("/", handlers.GetHouses)
	houses.Get("/:id", handlers.GetHouse)
	houses.Get("/:id/similar", handlers.GetSimilarHouses)
}

func SetupProtectedRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Get("/me", handlers.GetCurrentUser)
	auth.Put("/me", handlers.UpdateProfile)

	houses := api.Group("/houses")
	houses.Post("/", handlers.CreateHouse)
	houses.Put("/:id", handlers.UpdateHouse)
	houses.Delete("/:id", handlers.DeleteHouse)
	houses.Post("/:id/images", handlers.UploadHouseImages)
	houses.Put("/:id/status", handlers.UpdateHouseStatus)

	api.Get("/my-houses", handlers.GetMyHouses)

	favorites := api.Group("/favorites")
	favorites.Get("/", handlers.GetFavorites)
	favorites.Post("/:houseId", handlers.AddFavorite)
	favorites.Delete("/:houseId", handlers.RemoveFavorite)

	messages := api.Group("/messages")
	messages.Get("/", handlers.GetMessages)
	messages.Get("/sent", handlers.GetSentMessages)
	messages.Post("/", handlers.SendMessage)
	messages.Post("/:id/reply", handlers.ReplyMessage)
	messages.Put("/:id/read", handlers.MarkMessageAsRead)
}
