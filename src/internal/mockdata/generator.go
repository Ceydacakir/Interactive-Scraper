package mockdata

import (
	"fmt"
	"math/rand"
	"time"

	"interactive-scraper/src/internal/models"
)

var (
	vendors = []string{"SilkRoad", "AlphaBay", "Hydra", "DarkMarket", "WhiteHouseMarket", "TorRezac"}
	actions = []string{"Sale", "Leak", "Database", "Access", "Exploit", "Credit Cards", "SSN List"}
	targets = []string{"Bank", "Government", "Social Media", "Retail", "Crypto Exchange", "Email Provider"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateMockSource() models.Source {
	vendor := vendors[rand.Intn(len(vendors))]
	return models.Source{
		Name:             vendor,
		URL:              fmt.Sprintf("http://%s.onion", vendor),
		CriticalityScore: rand.Intn(10) + 1,
	}
}

func GenerateMockContent(sourceID uint) models.Content {
	action := actions[rand.Intn(len(actions))]
	target := targets[rand.Intn(len(targets))]

	title := fmt.Sprintf("%s: %s - %s", action, target, randomString(5))

	categories := []string{"Ransomware", "Database Leak", "Fraud", "Hacking Tool", "Insider Threat", "Drugs", "Weapons"}
	category := categories[rand.Intn(len(categories))]

	return models.Content{
		SourceID:    sourceID,
		Category:    category,
		Title:       title,
		RawContent:  fmt.Sprintf("Selling access to %s. Price: %d BTC. Contact: %s", target, rand.Intn(10), randomString(10)),
		PublishDate: time.Now(),
	}
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
