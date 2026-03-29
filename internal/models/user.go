package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         *string   `gorm:"uniqueIndex" json:"username"`
	FirstName        *string   `json:"first_name"`
	LastName         *string   `json:"last_name"`
	Email            string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash     *string   `json:"password_hash"`
	AuthProviderType string    `json:"auth_provider"` // e.g., "google", "apple", "facebook"
	AuthToken        string    `json:"auth_token"`    // Unique ID or token provided by the AuthProvider
	RoleID           uint      `gorm:"not null" json:"role_id"`
	DOB              *string   `json:"dob"`
	PhoneNumber      *string   `json:"phone_number"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	KlaviyoID        *string   `json:"klaviyo_id"`

	Roles   []*Role `gorm:"many2many:user_roles;"`
	IsBlock bool    `gorm:"default:false" json:"is_block"`
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`

	User *User
	Role *Role
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`

	Role       *Role
	Permission *Permission
}
