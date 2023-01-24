package router

import (
	"github.com/dp3why/mongofiber/common"
	"github.com/dp3why/mongofiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func getBooks(c *fiber.Ctx) error {
	coll := common.GetDBCollection("books")

	// find all books
	books := make([]models.Book, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		book := models.Book{}
		err := cursor.Decode(&book)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		books = append(books, book)
	}

	return c.Status(200).JSON(fiber.Map{"data": books})
}

func getBook(c *fiber.Ctx) error {
	coll := common.GetDBCollection("books")

	// find the book
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	book := models.Book{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": book})
}
