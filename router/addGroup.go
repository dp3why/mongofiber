package router

import (
	"github.com/gofiber/fiber/v2"
)

func AddBookGroup(app *fiber.App) {
	bookGroup := app.Group("/books")

	bookGroup.Get("/", getBooks)
	bookGroup.Get("/:id", getBook)
	bookGroup.Post("/", createBook)
	bookGroup.Put("/:id", updateBook)
	bookGroup.Delete("/:id", deleteBook)
}
