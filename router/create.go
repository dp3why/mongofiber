package router

import (
	"context"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/dp3why/mongofiber/common"
	"github.com/gofiber/fiber/v2"
	"github.com/pborman/uuid"
	"go.mongodb.org/mongo-driver/bson"
)


var BUCKET_NAME string =  os.Getenv("BUCKET_NAME")

// type createDTO struct {
// 	Id  string `json:"id" bson:"_id"`
// 	Title  string `json:"title" bson:"title"`
// 	Author string `json:"author" bson:"author"`
// 	Year   string `json:"year" bson:"year"`
// }

func createBook(c *fiber.Ctx) error {

    // Get the image from the form
    file, err := c.FormFile("image")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Failed to get image from form",
        })
    }
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

    // ===== Google GCS:  create a client for the Google Cloud Storage API =====
    client, err := storage.NewClient(c.Context())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to create client for Google Cloud Storage",
        })
    }    

	//======= create / get a bucket =========== 
	bucket := client.Bucket(BUCKET_NAME)

    // create a new file object
    obj := bucket.Object(file.Filename)

    // upload the file to the bucket
    w := obj.NewWriter(c.Context())
    if _, err = io.Copy(w, src); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to upload image to Google Cloud Storage",
        })
    }

    if err := w.Close(); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to close writer",
        })
    }

	ctx := context.Background()

 	attrs, _ := obj.Attrs(ctx)
	

	id := uuid.New()

    //  ==== Get form values =========
    title := c.FormValue("title")
    author := c.FormValue("author")
    year := c.FormValue("year")

	b := bson.M{
		"_id": id,
        "title": title,
        "author": author,
        "year": year,
		"Url": attrs.MediaLink,
    } 


	// ===== mongo: create the book ============
	coll := common.GetDBCollection("books")
	result, err := coll.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create book",
			"message": err.Error(),
		})
	}

	// return the book
	return c.Status(201).JSON(fiber.Map{
		"result": result,
		"Url": attrs.MediaLink,
	})
}
