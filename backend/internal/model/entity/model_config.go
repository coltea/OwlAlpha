package entity

import "time"

type ModelConfig struct {
	ID        uint       `gorm:"primaryKey"`
	BaseURL   string     `gorm:"column:base_url;size:255;not null"`
	APIKey    string     `gorm:"column:api_key;size:512;not null"`
	Model     string     `gorm:"size:128;not null"`
	CheckedAt *time.Time `gorm:"column:checked_at"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
}

func (ModelConfig) TableName() string {
	return "model_configs"
}
