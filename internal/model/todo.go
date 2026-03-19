package model

import "time"

type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"-"`
	Title     string    `gorm:"size:64;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Views     uint      `gorm:"default:0;not null" json:"views"`
	Status    bool      `gorm:"not null;default:false" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	StartAt   time.Time `gorm:"not null" json:"start_at"`
	EndAt     time.Time `gorm:"not null" json:"end_at"`
}

type TodoQueryConditions struct {
	UserID   uint   `json:"-"`
	Page     int    `json:"page"`
	PageSize int    `json:"pagesize"`
	Status   *bool  `json:"status,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
}
