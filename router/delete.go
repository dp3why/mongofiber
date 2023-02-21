package router

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/dp3why/mongofiber/common"
	"github.com/dp3why/mongofiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var book models.Book


func deleteBook(c *fiber.Ctx) error {
	// get the id
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

	// Find the book
	coll := common.GetDBCollection("books")

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to find book",
			"message": err.Error(),
		})
	}

	// delete the book from MongoDB
	result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete book",
			"message": err.Error(),
		})
	}
	

	// Delete on GCS
	deleteFile(BUCKET_NAME, book.Filename)
	
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}


// Google source code: deleteFile removes specified object.
func deleteFile( bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	attrs, err := o.Attrs(ctx)
	if err != nil {
			return fmt.Errorf("object.Attrs: %v", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}

	fmt.Println("successfully deleted the image on GCS")

	return nil
}