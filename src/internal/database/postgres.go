package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"interactive-scraper/src/internal/models"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Istanbul",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			log.Println("Connected to DB")
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatal("DB Connection failed")
}

func Migrate() {
	DB.AutoMigrate(&models.Source{}, &models.Content{})
	log.Println("Database Migration Completed")
}

func Seed() {
	log.Println("Seeding database...")
	sources := []models.Source{
		{Name: "Dready Forum", URL: "http://dreadytofatroptsdj6io7l3xptbet6onoyno2yv7jicoxknyazubrad.onion", CriticalityScore: 8},
		{Name: "Ramble", URL: "http://rambleeeqrhty6s5jgefdfdtc6tfgg4jj6svr4jpgk4wjtg3qshwbaad.onion", CriticalityScore: 5},
		{Name: "Darko", URL: "http://darkobds5j7xpsncsexzwhzaotyc4sshuiby3wtxslq5jy2mhrulnzad.onion/", CriticalityScore: 7},
	}
	for _, s := range sources {
		var count int64
		DB.Model(&models.Source{}).Where("url = ?", s.URL).Count(&count)
		if count == 0 {
			log.Printf("Adding source: %s", s.Name)
			DB.Create(&s)
		} else {
			log.Printf("Source exists: %s", s.Name)
		}
	}
}
