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
	time.Sleep(10 * time.Second)
	database.Connect()
	database.Migrate()
	database.Seed()
	for {
		scrape()
		time.Sleep(60 * time.Second)
	}
}

func getClient() *http.Client {
	p := os.Getenv("TOR_PROXY")
	if p == "" {
		p = "socks5://tor:9050"
	}
	u, _ := url.Parse(p)
	d, _ := proxy.FromURL(u, proxy.Direct)
	return &http.Client{Transport: &http.Transport{Dial: d.Dial}, Timeout: 30 * time.Second}
}

func check(title, body string) string {
	t := strings.ToLower(title + " " + body)
	if strings.Contains(t, "ransom") || strings.Contains(t, "lockbit") {
		return "Ransomware"
	}
	if strings.Contains(t, "database") || strings.Contains(t, "sql") {
		return "Database Leak"
	}
	if strings.Contains(t, "card") || strings.Contains(t, "market") {
		return "Illegal Market"
	}
	if strings.Contains(t, "forum") {
		return "Forum"
	}
	return "General"
}

func scrape() {
	var sources []models.Source
	database.DB.Find(&sources)
	client := getClient()

	for _, s := range sources {
		resp, err := client.Get(s.URL)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		title := strings.TrimSpace(doc.Find("title").Text())
		body := strings.Join(strings.Fields(doc.Find("body").Text()), " ")
		if len(body) > 500 {
			body = body[:500]
		}

		database.DB.Create(&models.Content{
			SourceID:    s.ID,
			Category:    check(title, body),
			Title:       title,
			RawContent:  body,
			PublishDate: time.Now(),
		})
		log.Println("Saved:", s.Name)
	}
}
