package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	Comment struct {
		Id     int    `json:"id"`
		PostId int    `json:"post_id"`
		Text   string `json:"text"`
	}
)

func main() {
	db, _ := initDB()
	initApp(db)
}

func initDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=comments_ms port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(Comment{})

	return db, nil
}

func initApp(db *gorm.DB) {
	app := fiber.New()

	app.Use(cors.New())

	app.Post("/api/comments", func(c *fiber.Ctx) error {
		var comment Comment

		if err := c.BodyParser(&comment); err != nil {
			panic(err)
		}

		db.Create(&comment)

		return c.JSON(comment)
	})

	app.Get("/api/posts/:id/comments", func(c *fiber.Ctx) error {
		var comments []Comment

		db.Find(&comments, "post_id=?", c.Params("id"))

		return c.JSON(comments)
	})

	app.Listen(":3020")
}
