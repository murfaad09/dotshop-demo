package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Supplier struct {
	gorm.Model
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        sql.NullString `json:"name"`
	Address     sql.NullString `json:"address"`
	ContactInfo sql.NullString `json:"contact_info"`
	Email       sql.NullString `json:"email"`
	Phone       sql.NullString `json:"phone"`
	Website     sql.NullString `json:"website"`
}
