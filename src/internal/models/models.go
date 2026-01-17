package models

import (
	"time"

	"gorm.io/gorm"
)

type Source struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Name             string         `json:"name"`
	URL              string         `json:"url"`
	CriticalityScore int            `json:"criticality_score"`
}

type Content struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	SourceID uint   `json:"source_id"`
	Source   Source `json:"source"`

	Title       string    `json:"title"`
	RawContent  string    `json:"raw_content"`
	PublishDate time.Time `json:"publish_date"`
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"unique" json:"username"`
	Password  string         `json:"-"`
}
