package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Post struct {
		Id          int       `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Comments    []Comment `json:"comments" gorm:"-" default:"[]"`
	}
)

func main() {
	db, _ := initDB()
	initApp(db)
}

func initDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=posts_ms port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(Post{})

	return db, nil
}

func initApp(db *gorm.DB) {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/posts", func(c *fiber.Ctx) error {
		var posts []Post

		db.Find(&posts)

		// This code could be improved in a better way
		for i, p := range posts {
			r, err := http.Get(fmt.Sprintf("http://localhost:3020/api/posts/%d/comments", p.Id))

			if err != nil {
				return err
			}

			var comments []Comment

			json.NewDecoder(r.Body).Decode(&comments)

			posts[i].Comments = comments
		}

		return c.JSON(posts)
	})

	app.Post("api/post", func(c *fiber.Ctx) error {
		var post Post

		err := c.BodyParser(&post)
		if err != nil {
			return err
		}

		db.Create(&post)

		return c.JSON(post)
	})

	app.Listen(":3010")
}
