package main

import (
	"log"
	"time"

	"interactive-scraper/src/internal/database"
	"interactive-scraper/src/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	time.Sleep(5 * time.Second)
	database.Connect()

	engine := html.New("./src/web/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Static("/static", "./src/web/static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		if c.Cookies("auth") == "" {
			return c.Redirect("/login")
		}
		return c.Render("dashboard", fiber.Map{
			"User": "Admin",
		})
	})

	api := app.Group("/api")
	api.Post("/login", func(c *fiber.Ctx) error {
		type LoginReq struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var req LoginReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad request"})
		}

		if req.Username == "admin" && req.Password == "admin" {
			c.Cookie(&fiber.Cookie{
				Name:  "auth",
				Value: "token_123",
			})
			return c.JSON(fiber.Map{"success": true})
		}
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	})

	api.Get("/stats", func(c *fiber.Ctx) error {
		var totalContent int64
		var totalSources int64

		database.DB.Model(&models.Content{}).Count(&totalContent)
		database.DB.Model(&models.Source{}).Count(&totalSources)

		type CriticalityGroup struct {
			Score int
			Count int
		}
		var groups []CriticalityGroup

		database.DB.Model(&models.Source{}).Select("criticality_score as score, count(*) as count").Group("criticality_score").Scan(&groups)

		return c.JSON(fiber.Map{
			"total_content": totalContent,
			"total_sources": totalSources,
			"criticality":   groups,
		})
	})

	api.Get("/content", func(c *fiber.Ctx) error {
		var contents []models.Content
		database.DB.Preload("Source").Order("created_at desc").Limit(50).Find(&contents)
		return c.JSON(contents)
	})

	log.Fatal(app.Listen(":3000"))
}
