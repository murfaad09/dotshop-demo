package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ReviewID uint   `gorm:"not null" json:"review_id"`
	UserID   uint   `gorm:"not null" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserID" json:"user"`
	Content  string `gorm:"type:text;not null" json:"content"`
}
