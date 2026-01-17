package main

import (
	"log"
	"math/rand"
	"time"

	"interactive-scraper/src/internal/database"
	"interactive-scraper/src/internal/mockdata"
	"interactive-scraper/src/internal/models"
)

func main() {
	log.Println("Starting Scraper Service (Mock)...")

	// Wait for DB to be ready (rudimentary check, retry loop)
	time.Sleep(5 * time.Second) // Give Postgres some time to start in Docker
	database.Connect()

	// Seed some sources if empty
	seedSources()

	// Main Loop
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			scrapeCycle()
		}
	}
}

func seedSources() {
	var count int64
	database.DB.Model(&models.Source{}).Count(&count)
	if count == 0 {
		log.Println("Seeding sources...")
		for i := 0; i < 5; i++ {
			source := mockdata.GenerateMockSource()
			database.DB.Create(&source)
		}
	}
}

func scrapeCycle() {
	log.Println("Running scrape cycle...")

	// Get a random source
	var source models.Source
	result := database.DB.Order("RANDOM()").First(&source)
	if result.Error != nil {
		log.Println("No sources found to scrape.")
		return
	}

	// Generate 1-3 new posts
	numPosts := rand.Intn(3) + 1
	for i := 0; i < numPosts; i++ {
		content := mockdata.GenerateMockContent(source.ID)
		if err := database.DB.Create(&content).Error; err != nil {
			log.Printf("Failed to save content: %v", err)
		} else {
			log.Printf("Saved new content: %s", content.Title)
		}
	}
}
