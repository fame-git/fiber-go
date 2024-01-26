package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// Book struct
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func checkMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(), c.Method(), start,
	)

	return c.Next()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	books = append(books, Book{ID: 1, Title: "Fail", Author: "Fail"})
	books = append(books, Book{ID: 2, Title: "FF", Author: "FF"})

	app.Use(checkMiddleware)

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)
	app.Get("/test-html", testHTML)

	app.Get("/config", getEnv)

	app.Listen(":8080")
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/"+file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File upload complete!")
}

func testHTML(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello World",
		"Name":  "Fail",
	})
}

func getEnv(c *fiber.Ctx) error {
	secret := os.Getenv("SECRET")

	if secret == "" {
		secret = "defualtsecret"
	}

	return c.JSON(fiber.Map{
		"SECRET": os.Getenv("SECRET"),
	})
}
