package router

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/dp3why/mongofiber/common"
	"github.com/gofiber/fiber/v2"
	"github.com/pborman/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// type createDTO struct {
// 	Id  string `json:"id" bson:"_id"`
// 	Title  string `json:"title" bson:"title"`
// 	Author string `json:"author" bson:"author"`
// 	Year   string `json:"year" bson:"year"`
// }

var BUCKET_NAME string =  os.Getenv("BUCKET_NAME")

func saveToGCS(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
    // Open the uploaded file
    src, err := file.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open uploaded file: %w", err)
    }
    defer src.Close()

    // Create a client for the Google Cloud Storage API
    client, err := storage.NewClient(c.Context())
    if err != nil {
        return "", fmt.Errorf("failed to create client for Google Cloud Storage: %w", err)
    }

    // Create/get a bucket
    bucket := client.Bucket(BUCKET_NAME)

    // Create a new file object
    obj := bucket.Object(file.Filename)

    // Upload the file to the bucket
    w := obj.NewWriter(c.Context())
    if _, err = io.Copy(w, src); err != nil {
        return "", fmt.Errorf("failed to upload image to Google Cloud Storage: %w", err)
    }

    if err := w.Close(); err != nil {
        return "", fmt.Errorf("failed to close writer: %w", err)
    }

	if err := obj.ACL().Set(c.Context(), storage.AllUsers, storage.RoleReader); err != nil {
		panic(err)
 	}
    // Get the URL of the uploaded file
    ctx := context.Background()
    attrs, _ := obj.Attrs(ctx)
    return attrs.MediaLink, nil
}



// =========== Create ===============
func createBook(c *fiber.Ctx) error {
    // Get the image from the form
    file, err := c.FormFile("image")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "failed to get image from form",
        })
    }

    // Upload the image to GCS
    url, err := saveToGCS(c, file)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Get form values
    title := c.FormValue("title")
    author := c.FormValue("author")
    year := c.FormValue("year")

    // Create the book
    coll := common.GetDBCollection("books")
    _, err = coll.InsertOne(c.Context(), bson.M{
        "_id":  uuid.New(),
        "title": title,
        "author": author,
        "year": year,
        "Url": url,
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error":   "Failed to create book",
            "message": err.Error(),
        })
    }

    // Return the book
    return c.Status(201).JSON(fiber.Map{
        "result": "Book created successfully",
        "Url": url,
    })
}

