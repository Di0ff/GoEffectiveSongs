package models

import (
	"gorm.io/gorm"
	"time"
)

type Songs struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	Group     string         `gorm:"index" json:"group"`
	Song      string         `gorm:"index" json:"song"`
	Release   time.Time      `json:"release_date"`
	Link      string         `gorm:"unique" json:"link"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Verses struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	SongID    int64     `gorm:"index" json:"song_id"`
	Number    int       `json:"number"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
