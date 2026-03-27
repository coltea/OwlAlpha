package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:64;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	Role      string    `gorm:"size:32;not null;default:admin"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
