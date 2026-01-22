package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"interactive-scraper/src/internal/database"
	"interactive-scraper/src/internal/models"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
)

func main() {
	log.Println("Starting Real Tor Scraper Service...")

	// Wait for DB to be ready
	time.Sleep(10 * time.Second)
	database.Connect()

	// Seed some initial real onion sites if empty
	seedRealSources()

	// Main Loop
	ticker := time.NewTicker(60 * time.Second) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			scrapeCycle()
		}
	}
}

func seedRealSources() {
	var count int64
	database.DB.Model(&models.Source{}).Count(&count)
	if count == 0 {
		log.Println("Seeding initial real sources...")
		sources := []models.Source{
			{Name: "Dready Forum", URL: "http://dreadytofatroptsdj6io7l3xptbet6onoyno2yv7jicoxknyazubrad.onion", CriticalityScore: 8},
			{Name: "Ramble", URL: "http://rambleeeqrhty6s5jgefdfdtc6tfgg4jj6svr4jpgk4wjtg3qshwbaad.onion", CriticalityScore: 5},
			{Name: "BFD Forum", URL: "http://bfdforumon7c2iprvgeqmdlbczvwahbqgz2y7ft5uodmijfl4tbqvnad.onion", CriticalityScore: 6},
		}
		for _, s := range sources {
			database.DB.Create(&s)
		}
	}
}

func getTorClient() *http.Client {
	torProxy := os.Getenv("TOR_PROXY") // e.g. "socks5://tor:9050"
	if torProxy == "" {
		torProxy = "socks5://tor:9050" // Default in docker
	}

	proxyURL, err := url.Parse(torProxy)
	if err != nil {
		log.Fatalf("Failed to parse proxy URL: %v", err)
	}

	// Create a SOCKS5 dialer
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to create SOCKS5 dialer: %v", err)
	}

	tr := &http.Transport{
		Dial: dialer.Dial,
	}

	return &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
}

func scrapeCycle() {
	log.Println("Starting scrape cycle via Tor...")

	var sources []models.Source
	database.DB.Find(&sources)

	client := getTorClient()

	for _, source := range sources {
		log.Printf("Scraping %s (%s)...", source.Name, source.URL)

		resp, err := client.Get(source.URL)
		if err != nil {
			log.Printf("Failed to fetch %s: %v", source.URL, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("Status code %d for %s", resp.StatusCode, source.URL)
			continue
		}

		// Parse HTML
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Failed to parse HTML for %s: %v", source.URL, err)
			continue
		}

		title := doc.Find("title").Text()
		title = strings.TrimSpace(title)
		if title == "" {
			title = "No Title Found"
		}

		// Basic content extraction (just grabbing details to save space)
		// In a real scenario, you'd want more sophisticated extraction
		bodyText := doc.Find("body").Text()
		bodyText = strings.Join(strings.Fields(bodyText), " ") // Clean usage of whitespace
		if len(bodyText) > 500 {
			bodyText = bodyText[:500] + "..."
		}

		// Save to DB
		content := models.Content{
			SourceID:    source.ID,
			Category:    "General", // We can't easily categorize without analysis, default to General
			Title:       title,
			RawContent:  bodyText,
			PublishDate: time.Now(),
		}

		if err := database.DB.Create(&content).Error; err != nil {
			log.Printf("Failed to save content for %s: %v", source.Name, err)
		} else {
			log.Printf("Successfully scraped and saved: %s", source.Name)
		}
	}
	log.Println("Scrape cycle completed.")
}
