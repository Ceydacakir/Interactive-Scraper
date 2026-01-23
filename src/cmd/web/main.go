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
	database.Migrate()
	database.Seed()

	engine := html.New("./src/web/templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Static("/static", "./src/web/static")

	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/login") })
	app.Get("/login", func(c *fiber.Ctx) error { return c.Render("login", fiber.Map{}) })

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		if c.Cookies("auth") == "" {
			return c.Redirect("/login")
		}
		return c.Render("dashboard", fiber.Map{"User": "Admin"})
	})

	api := app.Group("/api")

	api.Post("/login", func(c *fiber.Ctx) error {
		type Req struct {
			User string `json:"username"`
			Pass string `json:"password"`
		}
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.SendStatus(400)
		}
		if req.User == "admin" && req.Pass == "admin" {
			c.Cookie(&fiber.Cookie{Name: "auth", Value: "token_123"})
			return c.JSON(fiber.Map{"success": true})
		}
		return c.SendStatus(401)
	})

	api.Get("/stats", func(c *fiber.Ctx) error {
		var contentCount, sourceCount int64
		database.DB.Model(&models.Content{}).Count(&contentCount)
		database.DB.Model(&models.Source{}).Count(&sourceCount)

		var crit []struct{ Score, Count int }
		database.DB.Model(&models.Source{}).Select("criticality_score as score, count(*) as count").Group("criticality_score").Scan(&crit)

		var cats []struct {
			Category string
			Count    int
		}
		database.DB.Model(&models.Content{}).Select("category, count(*) as count").Group("category").Scan(&cats)

		return c.JSON(fiber.Map{
			"total_content": contentCount,
			"total_sources": sourceCount,
			"criticality":   crit,
			"categories":    cats,
		})
	})

	api.Get("/content", func(c *fiber.Ctx) error {
		var list []models.Content
		database.DB.Preload("Source").Order("created_at desc").Limit(50).Find(&list)
		return c.JSON(list)
	})

	api.Get("/sources", func(c *fiber.Ctx) error {
		var list []models.Source
		database.DB.Order("id asc").Find(&list)
		return c.JSON(list)
	})

	api.Post("/sources", func(c *fiber.Ctx) error {
		var req models.Source
		if err := c.BodyParser(&req); err != nil {
			return c.SendStatus(400)
		}
		req.CriticalityScore = 5
		database.DB.Create(&req)
		return c.JSON(fiber.Map{"success": true, "id": req.ID})
	})

	api.Delete("/sources/:id", func(c *fiber.Ctx) error {
		database.DB.Delete(&models.Source{}, c.Params("id"))
		return c.JSON(fiber.Map{"success": true})
	})

	api.Post("/sources/:id/criticality", func(c *fiber.Ctx) error {
		var req struct {
			Score int `json:"score"`
		}
		c.BodyParser(&req)
		var s models.Source
		database.DB.First(&s, c.Params("id"))
		s.CriticalityScore = req.Score
		database.DB.Save(&s)
		return c.JSON(fiber.Map{"success": true})
	})

	log.Fatal(app.Listen(":3000"))
}
