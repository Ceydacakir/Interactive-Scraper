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
	seed()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			scrape()
		}
	}
}

func seed() {
	var count int64
	database.DB.Model(&models.Source{}).Count(&count)
	if count == 0 {
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

func getClient() *http.Client {
	p := os.Getenv("TOR_PROXY")
	if p == "" {
		p = "socks5://tor:9050"
	}
	u, _ := url.Parse(p)
	d, _ := proxy.FromURL(u, proxy.Direct)
	return &http.Client{
		Transport: &http.Transport{Dial: d.Dial},
		Timeout:   30 * time.Second,
	}
}

func check(t, b string) string {
	txt := strings.ToLower(t + " " + b)
	if strings.Contains(txt, "ransom") || strings.Contains(txt, "lockbit") {
		return "Ransomware"
	}
	if strings.Contains(txt, "database") || strings.Contains(txt, "sql") {
		return "Database Leak"
	}
	if strings.Contains(txt, "card") || strings.Contains(txt, "market") {
		return "Illegal Market"
	}
	if strings.Contains(txt, "forum") {
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
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			continue
		}

		title := strings.TrimSpace(doc.Find("title").Text())
		body := strings.Join(strings.Fields(doc.Find("body").Text()), " ")
		if len(body) > 500 {
			body = body[:500]
		}

		cat := check(title, body)

		c := models.Content{
			SourceID:    s.ID,
			Category:    cat,
			Title:       title,
			RawContent:  body,
			PublishDate: time.Now(),
		}
		database.DB.Create(&c)
		log.Println("Saved:", s.Name)
	}
}
