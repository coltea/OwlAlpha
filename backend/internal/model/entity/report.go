package entity

import "time"

type Report struct {
	ID             uint      `gorm:"primaryKey"`
	TradeDate      string    `gorm:"size:16;index;not null"`
	StockCode      string    `gorm:"size:16;index;not null"`
	StockName      string    `gorm:"size:64;not null"`
	Summary        string    `gorm:"type:text;not null"`
	RiskLevel      string    `gorm:"size:32;not null"`
	Recommendation string    `gorm:"size:32;not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (Report) TableName() string {
	return "reports"
}
