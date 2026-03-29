package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/harishash/dotshop-be/internal/constants"
	"github.com/harishash/dotshop-be/internal/dto"
	domain "github.com/harishash/dotshop-be/internal/models"
	user_service "github.com/harishash/dotshop-be/internal/services"
	"github.com/harishash/dotshop-be/internal/utils/auth"
	"github.com/harishash/dotshop-be/internal/utils/errors"
)

type UserHandler struct {
	userService user_service.IUserService
}

var validate *validator.Validate

func NewUserHandler(userService user_service.IUserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Call service to get users
	users, err := u.userService.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	// Return users as JSON response
	return c.JSON(users)
}

// GetUserByID
//
//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Security		BearerAuth
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	domain.User
//	@Router			/users/{id} [get]
func (u *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := u.userService.GetUserByID(uint(id))
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Failed to get user"})
	}

	return c.JSON(user)
}

// CreateUser
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			User
//	@Accept			json
//	@Param			body	body	dto.SingupUserRequest	true	"Signup User Request"
//
//	@Produce		json
//	@Body			user  dto.SingupUserRequest
//	@Success		200	{object}	domain.User
//	@Router			/signup [post]
func (u *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parse request body to User struct
	var body dto.SingupUserRequest
	var invalidInputError constants.Errors = constants.InvalidInput
	var userExistsError constants.Errors = constants.UserAlreadyExists
	validate = validator.New()
	validate.RegisterValidation("password", passwordValidator)
	var user domain.User

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": invalidInputError.String(),
		})
	}

	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			fieldName := err.Field()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": customValidationMessages[fieldName],
			})
		}
	}

	role, _ := u.userService.RoleExists(body.RoleID)
	if !role {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": invalidInputError.String(),
		})
	}

	passwordHash := hashPassword(*body.Password)
	user = domain.User{
		Email: body.Email,
		// Username:     &username,
		FirstName:    &body.FirstName,
		LastName:     &body.LastName,
		PasswordHash: &passwordHash,
		RoleID:       body.RoleID,
		CreatedAt:    time.Now(),
	}

	newUser, err := u.userService.CreateUser(user)
	if err != nil {
		if err.Error() == userExistsError.String() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": userExistsError.String(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user, error: " + err.Error(),
		})
	}

	if newUser.IsBlock {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "user is blocked",
		})
	}

	tokenString, err := auth.GenerateJWT(body.Email, uint64(body.RoleID), uint64(newUser.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.FailedTokenCreation.String(),
		})
	}

	// Return created user as JSON response
	return c.Status(fiber.StatusCreated).JSON(constants.LoginResponse{
		AccessToken: tokenString,
		UserId:      int(newUser.ID),
		Email:       newUser.Email,
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		RoleId:      newUser.RoleID,
	})
}

// SignIn
//
//	@Summary		Authorize user
//	@Description	Authorize user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SigninUserRequest	true	"User"
//	@Success		200		{object}	constants.LoginResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Router			/signin [post]
func (u *UserHandler) Authorize(c *fiber.Ctx) error {
	body := new(dto.SigninUserRequest)
	validate = validator.New()
	// validate.RegisterValidation("password", passwordValidator)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidInput.String(),
		})
	}

	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, curatorId, err := u.userService.ValidateUser(body.Email, hashPassword(body.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": constants.InvalidUser.String(),
		})
	}
	if user.IsBlock {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "account blocked",
		})
	}
	tokenString, err := auth.GenerateJWT(body.Email, uint64(user.RoleID), uint64(user.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.FailedTokenCreation.String(),
		})
	}

	if user.RoleID != 2 {
		curatorId = nil
	}

	return c.Status(fiber.StatusOK).JSON(constants.LoginResponse{
		AccessToken: tokenString,
		UserId:      int(user.ID),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		RoleId:      user.RoleID,
		CuratorId:   curatorId,
	})

}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[\d]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(password)
	)
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// GetProfile Get Profile
//
//	@Summary		This endpoint is used to get our profile
//	@Description	This endpoint is used to get our profile
//	@Security		BearerAuth
//	@Tags			User
//	@Accept			application/json
//	@Success		200	{object}	dto.ProfileResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/user/profile [get]
func (u *UserHandler) GetProfile(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	profile, err := u.userService.GetProfile(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

func (u *UserHandler) ValidateToken(c *fiber.Ctx) error {

	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token is required"})
	}

	middleware := auth.JWTMiddlewareValidation(token)
	var err error
	if !middleware {
		err = errors.New("invalid token") // Assign a value to the err variable
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

// UpdateUserProfile Update consumer profile
//
//	@Summary		This endpoint is used to update consumer profile
//	@Description	This endpoint is used to update consumer profile
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.UpdateConsumerProfileRequest	true	"Update Consumer Profile Request"
//
//	@Success		200		{object}	dto.UpdateConsumerProfileResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/consumer/profile [put]
func (u *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	body, err := parseBody[dto.UpdateConsumerProfileRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)
	user, err := u.userService.UpdateConsumerProfile(userId, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	tokenString, err := auth.GenerateJWT(body.Email, uint64(user.RoleID), uint64(user.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.FailedTokenCreation.String(),
		})
	}

	response := &dto.UpdateConsumerProfileResponse{
		AccessToken: tokenString,
		UserId:      uint64(user.ID),
		FirstName:   *user.FirstName,
		LastName:    *user.LastName,
		Email:       user.Email,
		DOB:         *user.DOB,
		PhoneNumber: *user.PhoneNumber,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// AddUserAddress Add User Address
//
//	@Summary		This endpoint is used to add consumer address
//	@Description	This endpoint is used to add consumer address
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.AddConsumerAddressRequest	true	"Update Consumer Address Request"
//
//	@Success		200		{object}	dto.AddConsumerAddressResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/consumer/address [post]
func (u *UserHandler) AddUserAddress(c *fiber.Ctx) error {
	body, err := parseBody[dto.AddConsumerAddressRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)
	resp, err := u.userService.AddConsumerAddress(userId, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetAllUserAddress Get All User Addresses
//
//	@Summary		This endpoint is used to get all user addresses
//	@Description	This endpoint is used to get all user addresses
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			application/json
//
//	@Success		200	{object}	[]dto.AddConsumerAddressResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/consumer/address [get]
func (u *UserHandler) GetAllUserAddress(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint64)
	resp, err := u.userService.GetAllUserAddresses(uint(userId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetUserAddressById Get User Address By Id
//
//	@Summary		This endpoint is used to get user address by id
//	@Description	This endpoint is used to get user address by address id
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id	path		string	true	"Address ID"
//
//	@Success		200	{object}	dto.AddConsumerAddressResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/consumer/address/{id} [get]
func (u *UserHandler) GetUserAddressById(c *fiber.Ctx) error {
	addressIdStr := c.Params("id")
	if addressIdStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address Id is required")
	}

	addressId, err := strconv.ParseUint(addressIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid cart ID")
	}

	userId := c.Locals("user_id").(uint64)
	resp, err := u.userService.GetUserAddressById(userId, addressId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateAddressById Update Consumer Address By Id
//
//	@Summary		This endpoint is used to update the address with id
//	@Description	This endpoint is used to update the address with address id
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id		path		string								true	"Address ID"
//	@Param			body	body		dto.UpdateConsumerAddressRequest	true	"Update Address Request"
//
//	@Success		200		{object}	dto.AddConsumerAddressResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/consumer/address/{id} [patch]
func (u *UserHandler) UpdateAddressById(c *fiber.Ctx) error {
	addressIdStr := c.Params("id")
	if addressIdStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address Id is required")
	}

	addressId, err := strconv.ParseUint(addressIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid cart ID")
	}

	body, err := parseBody[dto.UpdateConsumerAddressRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)
	address, err := u.userService.UpdateUserAddress(userId, addressId, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(address)
}

// DeleteAddress Delete Address
//
//	@Summary		This endpoint is used to delete a address by id
//	@Description	This endpoint is used to delete a address by address id
//	@Tags			User
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			id	path		string	true	"Address ID"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/consumer/address/{id} [delete]
func (u *UserHandler) DeleteAddress(c *fiber.Ctx) error {
	addressIdStr := c.Params("id")
	if addressIdStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Address Id is required")
	}

	addressId, err := strconv.ParseUint(addressIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid cart ID")
	}

	userId := c.Locals("user_id").(uint64)
	if err := u.userService.DeleteAddress(userId, addressId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Address deleted successfully",
	})
}

// OrdersList Get Consumer Order List
//
//	@Summary		Get Consumer Order List
//	@Description	This endpoint is used to get consumer order list
//	@Tags			User
//	@Accept			application/json
//	@Security		BearerAuth
//	@Success		200	{object}	[]dto.OrdersListResponse
//	@Failure		400	{object}	fiber.Error
//	@Failure		401	{object}	fiber.Error
//	@Failure		403	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Router			/consumer/orders [get]
func (h *UserHandler) ConsumerOrdersList(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint64)
	resp, err := h.userService.ConsumerOrdersList(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resp)
}

// SendForgotPasswordEmail Send Forgot Password Email
//
//	@Summary		Send Forgot Password Email
//	@Description	This endpoint is used to send email for forgot password
//	@Tags			User
//	@Accept			application/json
//	@Param			body	body		dto.SendEmailForgotPasswordRequest	true	"Forgot Password Request"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/forgot-password/send-email [post]
func (h *UserHandler) SendForgotPasswordEmail(c *fiber.Ctx) error {
	body, err := parseBody[dto.SendEmailForgotPasswordRequest](c)
	if err != nil {
		return err
	}

	if err := h.userService.SendForgotPasswordEmail(body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

// ForgotPassword Forgot Password
//
//	@Summary		This endpoint is used for forgot password
//	@Description	This endpoint is used for forgot password
//	@Tags			User
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			body	body		dto.ForgotPasswordRequest	true	"Change Password Request"
//
//	@Success		200		{object}	dto.UpdatePasswordResponse
//	@Failure		400		{object}	fiber.Error
//	@Router			/forgot-password [patch]
func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	body, err := parseBody[dto.ForgotPasswordRequest](c)
	if err != nil {
		return err
	}

	validate = validator.New()
	validate.RegisterValidation("password", passwordValidator)
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId := c.Locals("user_id").(uint64)
	res, err := h.userService.ForgotPassword(body, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
