package model

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	ACTIVE   Status = "active"
	REJECTED Status = "rejected"
	PENDING  Status = "pending"
	BLOCKED  Status = "blocked"
)

type SocialMediaType string

const (
	Facebook  SocialMediaType = "Facebook"
	Twitter   SocialMediaType = "Twitter"
	Instagram SocialMediaType = "Instagram"
	LinkedIn  SocialMediaType = "LinkedIn"
)

type Curator struct {
	ID        uint           `gorm:"type:bigint; not null; primaryKey" json:"id"`
	UserID    uint           `gorm:"type:bigint; not null; unique" json:"user_id"`
	Name      string         `gorm:"type:varchar(255); not null" json:"name"`
	ShopName  string         `gorm:"type:varchar(255); not null; uniqueIndex" json:"shop_name"`
	Bio       string         `gorm:"type:text" json:"bio"`
	CreatedAt time.Time      `gorm:"created_at; default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at; default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	NumberofFollowers uint   `gorm:"type:bigint; not null; default:0" json:"number_of_followers"`
	ProfileImageURL   string `gorm:"type:varchar(255)" json:"profile_image_url"`
	CoverImageURL     string `gorm:"type:varchar(255)" json:"cover_image_url"`
	Status            Status `gorm:"type:text; default:pending" json:"status"`
	IsBlock           bool   `gorm:"default:false" json:"is_block"`

	User             User
	CuratorProducts  []CuratorProduct   `gorm:"foreignKey:CuratorID"`
	Looks            []*Look            `gorm:"foreignKey:CuratorID" json:"looks"`
	Collections      []*Collection      `gorm:"foreignKey:CuratorID" json:"collections"`
	SocialMediaLinks []*SocialMediaLink `gorm:"foreignKey:CuratorID" json:"social_media_links"`
	BankInformation  []*BankInformation `gorm:"foreignKey:CuratorID" json:"bank_information"`
	Product          []*Product         `gorm:"many2many:curator_products;"`
	PayoutHistory    []*PayoutHistory   `gorm:"foreignKey:CuratorID" json:"payout_history"`
	Payout           *Payout            `gorm:"foreignKey:CuratorID" json:"payout"`
}

type SocialMediaLink struct {
	ID           uint            `gorm:"type:bigint; not null; primaryKey"`
	CuratorID    uint            `gorm:"type:bigint; not null"`
	Type         SocialMediaType `gorm:"type:varchar(20); not null"`
	URL          string          `gorm:"type:text;"`
	AccessToken  string          `gorm:"type:text; "`
	OpenID       string          `gorm:"type:text; "`
	RefreshToken string          `gorm:"type:text;"`
	CreatedAt    time.Time       `gorm:"created_at; default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time       `gorm:"updated_at; default:CURRENT_TIMESTAMP"`

	Curator *Curator `gorm:"foreignKey:CuratorID;references:ID"`
}

type BankInformation struct {
	ID          uint      `gorm:"type:bigint; not null; primaryKey"`
	CuratorID   uint      `gorm:"type:bigint; not null"`
	Location    string    `gorm:"type:varchar(255); not null"`
	FirstName   string    `gorm:"type:varchar(255); not null"`
	LastName    string    `gorm:"type:varchar(255); not null"`
	DateOfBirth string    `gorm:"type:varchar(255); not null"`
	CreatedAt   time.Time `gorm:"created_at; default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"updated_at; default:CURRENT_TIMESTAMP"`

	Curator *Curator `gorm:"foreignKey:CuratorID;references:ID"`
}

type AccountDetails struct {
	ID             uint      `gorm:"type:bigint; not null; primaryKey"`
	BankID         uint      `gorm:"type:bigint; not null"`
	BankAddress    string    `gorm:"type:varchar(255); not null"`
	BankName       string    `gorm:"type:varchar(255); not null"`
	BranchCode     string    `gorm:"type:varchar(255); not null"`
	AccountNumber  string    `gorm:"type:varchar(255); not null"`
	AccountName    string    `gorm:"type:varchar(255); not null"`
	AccountAddress string    `gorm:"type:varchar(255); not null"`
	IBAN           string    `gorm:"type:varchar(255); not null"`
	CreatedAt      time.Time `gorm:"created_at; default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"updated_at; default:CURRENT_TIMESTAMP"`

	Bank *BankInformation `gorm:"foreignKey:BankID;references:ID"`
}
