package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	model "github.com/harishash/dotshop-be/internal/models"
	"gorm.io/gorm"
)

// extractEmailFromContext extracts the email from the context locals
func extractEmailFromContext(c *fiber.Ctx) (string, error) {
	emailInterface := c.Locals("email")

	email, ok := emailInterface.(string)
	if !ok || email == "" {
		return "", fiber.ErrUnauthorized
	}

	return email, nil
}

// getUserByEmail retrieves the user from the database by email
func getUserByEmail(db *gorm.DB, email string) (model.User, error) {
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, fiber.ErrForbidden
		}
		return user, err
	}
	return user, nil
}

// RoleMiddleware checks if the user has the specified role
func RoleMiddleware(db *gorm.DB, roleNames []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email, err := extractEmailFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		user, err := getUserByEmail(db, email)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		var role model.Role
		if err := db.Where("id = ?", user.RoleID).First(&role).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		for _, roleName := range roleNames {
			if role.Name == roleName {
				return c.Next()
			}
		}

		fmt.Println("User does not have the specified role.")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}

// PermissionMiddleware checks if the user has the specified permission
func PermissionMiddleware(db *gorm.DB, permissionName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email, err := extractEmailFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		user, err := getUserByEmail(db, email)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Retrieve the permission by name
		var permission model.Permission
		if err := db.Where("name = ?", permissionName).First(&permission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Permission not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		// Check if the role has the required permission
		var rolePermission model.RolePermission
		if err := db.Where("role_id = ? AND permission_id = ?", user.RoleID, permission.ID).First(&rolePermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Forbidden",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		return c.Next()
	}
}
