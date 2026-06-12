package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/user-api/internal/handler"
)

// Register wires all application routes onto the given Fiber app.
func Register(app *fiber.App, userHandler *handler.UserHandler) {
	api := app.Group("/users")

	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.ListUsers)
	api.Get("/:id", userHandler.GetUser)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)
}
