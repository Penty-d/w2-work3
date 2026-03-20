package model

import "time"

type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"-"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` //查ai的，gorm建表对todo加入外键约束，让用户被删除时对应todo也被删
	Title     string    `gorm:"size:64;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Views     uint      `gorm:"default:0;not null" json:"view"`
	Status    bool      `gorm:"not null;default:false" json:"status"`
	CreatedAt time.Time `gorm:"not null" json:"create_at"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
}

type TodoQueryConditions struct {
	UserID   uint     `json:"-"`
	Page     int      `json:"page"`
	PageSize int      `json:"pagesize"`
	Status   *bool    `json:"status,omitempty"`
	Keywords []string `json:"keywords,omitempty"`
}
