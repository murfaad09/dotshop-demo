package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sync"
	"time"

	"github.com/harishash/dotshop-be/internal/config"
	domain "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Connection *gorm.DB
}

var dbinstance *Database
var dbonce sync.Once

func GetDatabaseConnection() *Database {
	dbonce.Do(func() {
		dbinstance = &Database{
			Connection: dbConnect(),
		}
	})
	return dbinstance
}

func dbConnect() *gorm.DB {
	c := config.GetConfig()

	db, err := gorm.Open(postgres.Open(c.DbUrl), &gorm.Config{})
	if err != nil {
		//log.Fatalf("Error connecting to database: %v", err) // demo mode 
		    log.Println("Running in demo mode (no DB)")

	}

	// err = db.AutoMigrate(
	// 	&domain.User{},
	// 	&domain.Role{},
	// 	&domain.Permission{},
	// 	&domain.UserRole{},
	// 	&domain.RolePermission{},
	// 	&domain.AccountDetails{},  // Ensure AccountDetails is created first
	// 	&domain.BankInformation{}, // Then create BankInformation
	// 	&domain.Curator{},
	// 	&domain.SocialMediaLink{},
	// 	&domain.Collection{},
	// 	&domain.Product{},
	// 	&domain.Cart{},
	// 	&domain.CartItem{},
	// 	&domain.CollectionProduct{},
	// 	&domain.Look{},
	// 	&domain.CuratorProduct{},
	// 	&domain.CollectionSection{},
	// 	&domain.CollectionSectionProduct{},
	// 	&domain.Variant{},
	// 	&domain.VariantOption{},
	// 	&domain.Category{},
	// 	&domain.SubCategory{},
	// 	&domain.AuditTrail{},
	// 	// &domain.Brand{},
	// 	&domain.VendorConfiguration{},
	// 	&domain.SupportTicket{},
	// 	&domain.Supplier{},
	// 	&domain.Order{},
	// 	&domain.ShippingInfo{},
	// 	&domain.Session{},
	// 	&domain.Review{},
	// 	&domain.Comment{},
	// 	&domain.Refund{},
	// 	&domain.Promotion{},
	// 	&domain.Payment{},
	// 	// &domain.OrderItem{},
	// 	&domain.Notification{},
	// 	&domain.GiftCard{},
	// 	// &domain.Favorite{},
	// 	&domain.Media{},
	// 	&domain.WishlistItem{},
	// 	&domain.Log{},
	// 	&domain.OrderVariants{},
	// 	&domain.Fulfillments{},
	// 	&domain.Payout{},
	// 	&domain.PayoutHistory{},
	// 	&domain.ReturnOrder{},
	// 	&domain.Notification{},
	// )
	// if err != nil {
	// 	log.Fatalf("Error migrating database: %v", err)
	// }

	// createDefaultRolesAndCurator(db)

	return db
}

func createDefaultRolesAndCurator(db *gorm.DB) {
	roles := []domain.Role{
		{Name: "user"},
		{Name: "curator"},
		{Name: "admin"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, domain.Role{Name: role.Name}).Error; err != nil {
			log.Fatalf("Error creating role %s: %v", role.Name, err)
		}
	}

	var curatorRole, adminRole domain.Role
	if err := db.Where("name = ?", "curator").First(&curatorRole).Error; err != nil {
		log.Fatalf("Error finding curator role: %v", err)
	}

	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Fatalf("Error finding admin role: %v", err)
	}

	curatorUser := domain.User{
		FirstName:    stringPointer("Dot"),
		LastName:     stringPointer("Shop"),
		Email:        config.GetConfig().DotShopStoreEmail,
		PasswordHash: stringPointer(hashPassword(config.GetConfig().DotShopStorePassword)),
		RoleID:       curatorRole.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	adminUser := domain.User{
		FirstName:    stringPointer("Admin"),
		LastName:     stringPointer("Admin"),
		Email:        config.GetConfig().DotShopAdminEmail,
		PasswordHash: stringPointer(hashPassword(config.GetConfig().DotShopAdminPassword)),
		RoleID:       adminRole.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	dotshopCurator := domain.Curator{
		Name:              "DotShop",
		ShopName:          "DotShop",
		Bio:               "DotShop",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		NumberofFollowers: 0,
		ProfileImageURL:   "defaultprofileurl",
		CoverImageURL:     "defaultcoverurl",
		Status:            domain.ACTIVE,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.FirstOrCreate(&curatorUser, domain.User{Email: curatorUser.Email}).Error; err != nil {
			return err
		}

		dotshopCurator.UserID = curatorUser.ID
		if err := tx.FirstOrCreate(&dotshopCurator, domain.Curator{UserID: dotshopCurator.UserID}).Error; err != nil {
			return err
		}

		if err := tx.FirstOrCreate(&adminUser, domain.User{Email: adminUser.Email}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error creating default curator: %v", err)
	}
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func stringPointer(s string) *string {
	return &s
}
