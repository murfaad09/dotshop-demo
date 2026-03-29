package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/dto"
	models "github.com/harishash/dotshop-be/internal/models"
	services "github.com/harishash/dotshop-be/internal/services"

	"github.com/harishash/dotshop-be/internal/utils/auth"

	"gorm.io/gorm"
)

var db *gorm.DB
var userService services.IUserService

// SetDependencies sets the database connection and user service in the handlers
func SetDependencies(database *gorm.DB, usrService services.IUserService) {
	db = database
	userService = usrService
}

// GoogleSSOCallback 	handles the Google SSO callback
//
//	@Summary		handles the Google SSO callback
//	@Description	handles the Google SSO callback
//	@Tags			GoogleSSO
//	@Accept			application/json
//	@Param			body	body		dto.GoogleCallbackRequest	true	"Google SSO"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/auth/google/callback [post]
func GoogleSSOCallback(c *fiber.Ctx) error {
	var googleResp dto.GoogleCallbackRequest

	if err := c.BodyParser(&googleResp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse body"})
	}

	// Verify Google access token
	tokenInfo, err := verifyGoogleToken(googleResp.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	var user models.User
	err = db.Where("email = ?", tokenInfo.Email).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to query user"})
	}
	if user.IsBlock {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "account blocked",
		})
	}
	var jwtToken string
	var userID int
	var roleID int

	if err == nil {
		// User exists, perform signin
		jwtToken, err = auth.GenerateJWT(user.Email, uint64(user.RoleID), uint64(user.ID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate JWT token")
		}
		userID = int(user.ID)
		roleID = int(user.RoleID)
	} else if err == gorm.ErrRecordNotFound {
		// User does not exist, create new user
		// Split the full name into parts
		parts := strings.SplitN(googleResp.Name, " ", 2)

		var firstName, lastName string

		if len(parts) == 2 {
			firstName = parts[0]
			lastName = parts[1]
		} else {
			firstName = parts[0]
			lastName = ""
		}

		newUser := models.User{
			Email:            tokenInfo.Email,
			Username:         &googleResp.Name,
			FirstName:        &firstName,
			LastName:         &lastName,
			RoleID:           constants.USER_ROLE_ID,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			AuthProviderType: constants.AuthProviderType,
			AuthToken:        googleResp.AccessToken,
		}
		user, err := userService.CreateUser(newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
		}

		jwtToken, err = auth.GenerateJWT(user.Email, uint64(user.RoleID), uint64(user.ID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate JWT token")
		}
		userID = int(user.ID)
		roleID = int(user.RoleID)
	}

	return c.JSON(fiber.Map{"token": jwtToken, "user_id": userID, "role_id": roleID})
}

func verifyGoogleToken(accessToken string) (dto.GoogleCallbackRequest, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=%s", accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return dto.GoogleCallbackRequest{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.GoogleCallbackRequest{}, fmt.Errorf("invalid token, status code: %d", resp.StatusCode)
	}

	var user dto.GoogleCallbackRequest
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return dto.GoogleCallbackRequest{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return user, nil
}
