package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserName     string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:60;not null" json:"-"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
}
