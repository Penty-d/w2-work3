package model

import "time"

type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"size:64;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Views     uint      `gorm:"default:0;not null" json:"views"`
	Status    bool      `gorm:"not null;default:false" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	StartAt   time.Time `gorm:"not null" json:"start_at"`
	EndAt     time.Time `gorm:"not null" json:"end_at"`
}
