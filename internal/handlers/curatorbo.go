package handlers

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	dto "github.com/harishash/dotshop-be/internal/dto"
	curatorboservices "github.com/harishash/dotshop-be/internal/services"

	"gorm.io/gorm"
)

type CuratorBOHandlers struct {
	service curatorboservices.CuratorBOService
}

func NewCuratorBOHandlers(service curatorboservices.CuratorBOService) *CuratorBOHandlers {
	return &CuratorBOHandlers{service: service}
}

// AddProduct creates a new product
//
//	@Summary		Create a new curator feature product
//	@Description	Create a new product
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CreateFeatureProductRequest	true	"Feature Product"
//	@Success		200		{object}	dto.CreateFeatureProductResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/addproduct [post]
func (h *CuratorBOHandlers) AddProduct(c *fiber.Ctx) error {
	var payload *dto.CreateFeatureProductRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	insertedProducts, err := h.service.AddProduct(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(insertedProducts)
}

// AddCollection creates a new collection
//
//	@Summary		Create a new curator collection
//	@Description	Create a new collection
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CreateCollectionRequest	true	"Collection Product"
//	@Success		200		{object}	dto.CreateCollectionResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/addcollection [post]
func (h *CuratorBOHandlers) AddCollection(c *fiber.Ctx) error {
	var payload *dto.CreateCollectionRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := payload.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	collection, err := h.service.AddCollection(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(collection)
}

// AddCollectionSection creates a new collection section
//
//	@Summary		Create a new collection section
//	@Description	Create a new collection section
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CreateCollectionSectionRequest	true	"Collection Section Product"
//	@Success		200		{object}	dto.CreateCollectionSectionResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/collection/addsection [post]
func (h *CuratorBOHandlers) AddCollectionSection(c *fiber.Ctx) error {
	var payload *dto.CreateCollectionSectionRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	section, err := h.service.AddCollectionSection(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(section)
}

// AddProductToSection adds a product to a section.
//
//	@Summary		Add a product to a section
//	@Description	Add a product to a section
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			section_id	path		string							true	"Section ID"
//	@Param			body		body		dto.AddProductToSectionRequest	true	"ProductS to add"
//	@Success		200			{object}	dto.AddProductToSectionResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/section/{section_id}/addproduct [post]
func (h *CuratorBOHandlers) AddProductToSection(c *fiber.Ctx) error {

	sectionID, err := strconv.ParseUint(c.Params("section_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid section ID"})
	}
	var payload *dto.AddProductToSectionRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	section, err := h.service.AddProductToSection(uint(sectionID), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(section)
}

// UpdateSection updates a section.
//
//	@Summary		Update section info
//	@Description	Update section info
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			section_id	path		string								true	"Section ID"
//	@Param			body		body		dto.UpdateCollectionSectionRequest	true	"Collection Section Product"
//	@Success		200			{object}	dto.UpdateCollectionSectionResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/section/{section_id} [put]
func (h *CuratorBOHandlers) UpdateSection(c *fiber.Ctx) error {

	sectionID, err := strconv.ParseUint(c.Params("section_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid section ID"})
	}

	var payload *dto.UpdateCollectionSectionRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}
	section, err := h.service.UpdateCollectionSection(payload, uint(sectionID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(section)
}

// DeleteProductFromSectionByID deletes a product from a section.
//
//	@Summary		Delete a product from a section
//	@Description	Delete a product from a section by its ID and the section's ID
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			section_id	path		string		true	"Section ID"
//	@Param			product_id	path		string		true	"Product ID"
//	@Success		200			{object}	fiber.Map	"Successfully deleted the product from the section"
//	@Failure		400			{object}	fiber.Error	"Bad Request"
//	@Failure		401			{object}	fiber.Error	"Unauthorized"
//	@Failure		403			{object}	fiber.Error	"Forbidden"
//	@Failure		404			{object}	fiber.Error	"Not Found"
//	@Router			/curator/section/{section_id}/product/{product_id} [delete]
func (h *CuratorBOHandlers) DeleteProductFromSectionByID(c *fiber.Ctx) error {

	sectionID, err := strconv.ParseUint(c.Params("section_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sectionID ID"})
	}
	productID := c.Params("product_id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Section ID is required"})
	}

	err = h.service.DeleteProductFromSectionByID(uint(sectionID), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}

// AddLook creates a new look
//
//	@Summary		Create a new curator look
//	@Description	Create a new look
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			body	body		dto.CreateLookRequest	true	"Look Product"
//	@Success		200		{object}	dto.CreateLookResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/addlook [post]
func (h *CuratorBOHandlers) AddLook(c *fiber.Ctx) error {
	var payload *dto.CreateLookRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := payload.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	look, err := h.service.AddLook(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(look)
}

// DeleteProductFromFeature deletes a product from a feature product
//
//	@Summary		Delete a product from a feature product
//	@Description	Delete a product from a feature product
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			curator_id	path		string	true	"Curator ID"
//	@Param			product_id	path		string	true	"Product ID"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/feature/product/{product_id} [delete]
func (h *CuratorBOHandlers) DeleteProductFromFeatureByID(c *fiber.Ctx) error {
	curatorIDStr := c.Params("curator_id")
	productID := c.Params("product_id")

	curatorID, err := strconv.ParseUint(curatorIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	err = h.service.DeleteFromFeatureProduct(uint(curatorID), productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "feature product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product from feature products"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}

// DeleteCollectionByID deletes a collection by collection ID
//
//	@Summary		Delete a collection by collection ID
//	@Description	Delete a collection by collection ID
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			collection_id	path		int	true	"Collection ID"
//	@Success		200				{object}	fiber.Map
//	@Failure		400				{object}	fiber.Error
//	@Failure		401				{object}	fiber.Error
//	@Failure		403				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collection/{collection_id} [delete]
func (h *CuratorBOHandlers) DeleteCollectionByID(c *fiber.Ctx) error {
	collectionIDStr := c.Params("collection_id")

	collectionID, err := strconv.ParseUint(collectionIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}
	err = h.service.DeleteCollectionByID(uint(collectionID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete collection"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Collection deleted successfully"})
}

// DeleteProductFromCollectionByID deletes a product from a collection by product ID
//
//	@Summary		Delete a product from a collection by product ID
//	@Description	Delete a product from a collection by product ID
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			collection_id	path		int		true	"Collection ID"
//	@Param			product_id		path		string	true	"Product ID"
//	@Success		200				{object}	fiber.Map
//	@Failure		400				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collection/{collection_id}/product/{product_id} [delete]
func (h *CuratorBOHandlers) DeleteProductFromCollectionByID(c *fiber.Ctx) error {
	collectionIDStr := c.Params("collection_id")
	productID := c.Params("product_id")

	collectionID, err := strconv.ParseUint(collectionIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid look ID"})
	}

	err = h.service.DeleteProductFromCollectionByID(uint(collectionID), productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product from collection"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}

// DeleteLookByID deletes a look by look ID
//
//	@Summary		Delete a look by look ID
//	@Description	Delete a look by look ID
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			look_id	path		int	true	"Look ID"
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Router			/curator/look/{look_id} [delete]
func (h *CuratorBOHandlers) DeleteLookByID(c *fiber.Ctx) error {
	lookIDStr := c.Params("look_id")

	lookID, err := strconv.ParseUint(lookIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid look ID"})
	}
	err = h.service.DeleteLookByID(uint(lookID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Look not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete look"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Look deleted successfully"})
}

// DeleteSectionByID deletes a section by ID
//
//	@Summary		Delete a section by ID
//	@Description	Delete a section by ID
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			section_id	path		int	true	"Section ID"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/section/{section_id} [delete]
func (h *CuratorBOHandlers) DeleteSectionByID(c *fiber.Ctx) error {
	sectionIDStr := c.Params("section_id")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sectionID ID"})
	}
	err = h.service.DeleteSectionByID(uint(sectionID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Section not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete look"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "section deleted successfully"})
}

// DeleteProductFromLookByID deletes a look product by look ID and product ID
//
//	@Summary		deletes a look product by look ID and product ID
//	@Description	deletes a look product by look ID and product ID
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			look_id		path		int		true	"Look ID"
//	@Param			product_id	path		string	true	"Product ID"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/look/{look_id}/product/{product_id} [delete]
func (h *CuratorBOHandlers) DeleteProductFromLookByID(c *fiber.Ctx) error {
	lookIDStr := c.Params("look_id")
	productID := c.Params("product_id")

	lookID, err := strconv.ParseUint(lookIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid look ID"})
	}

	err = h.service.DeleteProductFromLookByID(uint(lookID), productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Look not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product from look"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}

// func (h *CuratorBOHandlers) GetOrdersStatus(c *fiber.Ctx) error {
//     status, err := h.service.GetOrdersStatus()
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//     }
//     return c.JSON(status)
// }

// func (h *CuratorBOHandlers) GetPayoutDetails(c *fiber.Ctx) error {
//     details, err := h.service.GetPayoutDetails()
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//     }
//     return c.JSON(details)
// }

func (h *CuratorBOHandlers) Withdraw(c *fiber.Ctx) error {
	err := h.service.Withdraw(c.Body())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}

// func (h *CuratorBOHandlers) GetProfile(c *fiber.Ctx) error {
//     profile, err := h.service.GetProfile()
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//     }
//     return c.JSON(profile)
// }

// UpdateCuratorProfile Update curator profile
//
//	@Summary		This endpoint is used to update curator profile
//	@Description	This endpoint is used to update curator profile
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			body	body		dto.UpdateProfileRequest	true	"Update Curator Request"
//
//	@Success		200		{object}	dto.UpdateProfileResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/profile [put]
func (h *CuratorBOHandlers) UpdateProfile(c *fiber.Ctx) error {
	body, err := parseBody[dto.UpdateProfileRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)
	curator, user, err := h.service.UpdateProfile(userId, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := &dto.UpdateProfileResponse{
		UserId:          uint64(user.ID),
		CuratorId:       uint64(curator.ID),
		FirstName:       *user.FirstName,
		LastName:        *user.LastName,
		Bio:             curator.Bio,
		CoverImageURL:   curator.CoverImageURL,
		ProfileImageURL: curator.ProfileImageURL,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// InsertSocialMediaLink Insert Social Media Link
//
//	@Summary		This endpoint is used to insert social media link
//	@Description	This endpoint is used to insert social media link
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			curator_id	path		int							true	"Curator ID"
//	@Param			body		body		dto.SocialMediaLinksRequest	true	"Insert Social Media Link Request"
//
//	@Success		200			{object}	dto.CreateSocialMediaLinkResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/sociallink [post]
func (h *CuratorBOHandlers) InsertSocialMediaLink(c *fiber.Ctx) error {
	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return err
	}
	var request dto.SocialMediaLinksRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	var response interface{}

	exists, err := h.service.CheckSocialMediaLinkExists(request.Platform, uint(curatorID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if exists {
		response, err = h.service.UpdateSocialMediaLinks(&request, uint(curatorID))
	} else {
		response, err = h.service.AddSocialMediaLink(uint(curatorID), &request)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Swagger removed bcz this function not used
func (h *CuratorBOHandlers) UpdateSocialMediaLink(c *fiber.Ctx) error {

	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return err
	}

	linkID, err := strconv.ParseUint(c.Params("link_id"), 10, 64)
	if err != nil {
		return err
	}

	var links dto.SocialMediaLinksRequest
	if err := c.BodyParser(&links); err != nil {
		return err
	}

	response, err := h.service.
		UpdateSocialMediaLink(uint(curatorID), uint(linkID), &links)
	if err != nil {
		return err
	}

	return c.JSON(response)

}

// DeleteSocialMediaLink Delete Social Media Link
//
//	@Summary		This endpoint is used to delete social media link
//	@Description	This endpoint is used to delete social media link
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			curator_id	path		int	true	"Curator ID"
//	@Param			link_id		path		int	true	"Social Media Link ID"
//
//	@Success		200			{object}	dto.DeleteSocialMediaLinkResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/{curator_id}/sociallink/{link_id} [delete]
func (h *CuratorBOHandlers) DeleteSocialMediaLink(c *fiber.Ctx) error {
	linkID, err := strconv.ParseUint(c.Params("link_id"), 10, 64)
	if err != nil {
		return err
	}
	curatorID, err := strconv.ParseUint(c.Params("curator_id"), 10, 64)
	if err != nil {
		return err
	}
	response, err := h.service.RemoveSocialMediaLink(curatorID, linkID)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// ChangePassword Change Password
//
//	@Summary		This endpoint is used for change password from profile
//	@Description	This endpoint is used for change password from profile
//	@Tags			User
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			body	body		dto.UpdatePasswordRequest	true	"Change Password Request"
//
//	@Success		200		{object}	dto.UpdatePasswordResponse
//	@Failure		400		{object}	fiber.Error
//	@Router			/curator/profile/password [put]
func (h *CuratorBOHandlers) ChangePassword(c *fiber.Ctx) error {
	body, err := parseBody[dto.UpdatePasswordRequest](c)
	if err != nil {
		return err
	}
	validate = validator.New()
	validate.RegisterValidation("password", passwordValidator)

	if err := validate.Struct(body); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			fieldName := err.Field()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": customValidationMessages[fieldName],
			})
		}
	}

	err, res := h.service.ChangePassword(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetAllCurators Get all curators
//
//	@Summary		This endpoint is used for get all curators
//	@Description	This endpoint is used for get all curators with all the neccessary information
//	@Tags			CuratorBO
//	@Produce		application/json
//	@Success		200	{object}	[]dto.GetAllCuratorsResponse
//	@Failure		400	{object}	fiber.Error
//	@Router			/curator/all [get]
func (h *CuratorBOHandlers) GetAllCurators(c *fiber.Ctx) error {
	curators, err := h.service.GetAllCurators()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// var response dto.GetAllCuratorsResponse
	// for _, v := range curators {
	// 	response.Curators = append(response.Curators, dto.NewGetCuratorResponse(v))
	// }

	return c.Status(fiber.StatusOK).JSON(curators)
}

// GetCuratorWithCuratorID Get curator with curator id
//
//	@Summary		This endpoint is used for get curator with curator id
//	@Description	This endpoint is used for get curator with all the neccessary information
//	@Tags			CuratorBO
//	@Produce		application/json
//	@Param			curator_id	path		int	true	"Curator ID"
//	@Success		200			{object}	dto.GetCuratorResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id} [get]
func (h *CuratorBOHandlers) GetCuratorWithCuratorID(c *fiber.Ctx) error {
	curatorIDStr := c.Params("curator_id")
	if len(curatorIDStr) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "curator id not provided"})
	}

	curatorID, err := strconv.ParseUint(curatorIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	curator, err := h.service.GetCuratorByCuratorID(curatorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if curator.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "curator not found"})
	}
	return c.Status(fiber.StatusOK).JSON(dto.NewGetCuratorResponse(*curator))
}

// AddProductToFeature Add product to Feature
//
//	@Summary		This endpoint is used to add product to Feature
//	@Description	This endpoint is used to add product to Feature
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			feature_id	path		int								true	"Feature ID"
//	@Param			body		body		dto.AddProductToFeatureRequest	true	"Add Product to Feature Request"
//
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	fiber.Error
//	@Failure		401			{object}	fiber.Error
//	@Failure		403			{object}	fiber.Error
//	@Failure		404			{object}	fiber.Error
//	@Router			/curator/feature/{feature_id}/addproduct [post]
func (h *CuratorBOHandlers) AddProductToFeature(c *fiber.Ctx) error {
	featureIdStr := c.Params("feature_id")
	if len(featureIdStr) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "feature id not provided"})
	}

	featureId, err := strconv.ParseUint(featureIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid feature ID"})
	}

	body, err := parseBody[dto.AddProductToFeatureRequest](c)
	if err != nil {
		return err
	}
	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userId := c.Locals("user_id").(uint64)
	feature, err := h.service.AddProductToFeature(body, uint(featureId), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Products Added Successfully", "products": feature})
}

// AddProductToLook Add product to look
//
//	@Summary		This endpoint is used to add product to look
//	@Description	This endpoint is used to add product to look
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			look_id	path		int							true	"Look ID"
//	@Param			body	body		dto.AddProductToLookRequest	true	"Add Product to Look Request"
//
//	@Success		200		{object}	fiber.Map
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/look/{look_id}/addproduct [post]
func (h *CuratorBOHandlers) AddProductToLook(c *fiber.Ctx) error {
	lookIdStr := c.Params("look_id")
	if len(lookIdStr) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "look id not provided"})
	}

	lookId, err := strconv.ParseUint(lookIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid look ID"})
	}

	body, err := parseBody[dto.AddProductToLookRequest](c)
	if err != nil {
		return err
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userId := c.Locals("user_id").(uint64)
	look, err := h.service.AddProductToLook(body, uint(lookId), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Products Added Successfully", "products": look})
}

// AddProductToCollection Add product to collection
//
//	@Summary		This endpoint is used to add product to collection
//	@Description	This endpoint is used to add product to collection
//	@Tags			CuratorBO
//	@Accept			application/json
//	@Security		BearerAuth
//	@Param			collection_id	path		int									true	"Collection ID"
//	@Param			body			body		dto.AddProductToCollectionRequest	true	"Add Product to Collection Request"
//
//	@Success		200				{object}	fiber.Map
//	@Failure		400				{object}	fiber.Error
//	@Failure		401				{object}	fiber.Error
//	@Failure		403				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collection/{collection_id}/addproduct [post]
func (h *CuratorBOHandlers) AddProductToCollection(c *fiber.Ctx) error {
	collectionIdStr := c.Params("collection_id")
	if len(collectionIdStr) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "collection id not provided"})
	}

	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	body, err := parseBody[dto.AddProductToCollectionRequest](c)
	if err != nil {
		return err
	}

	if err := body.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userId := c.Locals("user_id").(uint64)
	collection, err := h.service.AddProductToCollection(body, uint(collectionId), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Products Added Successfully", "products": collection})
}

// UpdateCollection Update Collection
//
//	@Summary		This endpoint is used to update collection fields
//	@Description	This endpoint is used to update collection fields
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			collection_id	path		int							true	"Collection ID"
//	@Param			body			body		dto.UpdateCollectionRequest	true	"Update Collection Request"
//
//	@Success		200				{object}	dto.UpdateCollectionResponse
//	@Failure		400				{object}	fiber.Error
//	@Failure		401				{object}	fiber.Error
//	@Failure		403				{object}	fiber.Error
//	@Failure		404				{object}	fiber.Error
//	@Router			/curator/collection/{collection_id} [put]
func (h *CuratorBOHandlers) UpdateCollection(c *fiber.Ctx) error {
	collectionIdStr := c.Params("collection_id")
	if len(collectionIdStr) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "collection id not provided"})
	}

	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid collection ID"})
	}

	body, err := parseBody[dto.UpdateCollectionRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)

	resp, err := h.service.UpdateCollection(body, uint(collectionId), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateLook Update Look
//
//	@Summary		This endpoint is used to update look fields
//	@Description	This endpoint is used to update look fields
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Accept			application/json
//	@Param			look_id	path		int						true	"Look Id"
//	@Param			body	body		dto.UpdateLookRequest	true	"Update Look Request"
//
//	@Success		200		{object}	dto.UpdateLookResponse
//	@Failure		400		{object}	fiber.Error
//	@Failure		401		{object}	fiber.Error
//	@Failure		403		{object}	fiber.Error
//	@Failure		404		{object}	fiber.Error
//	@Router			/curator/look/{look_id} [put]
func (h *CuratorBOHandlers) UpdateLook(c *fiber.Ctx) error {
	lookIdStr := c.Params("look_id")
	if len(lookIdStr) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "look id not provided"})
	}

	lookId, err := strconv.ParseUint(lookIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid look Id"})
	}

	body, err := parseBody[dto.UpdateLookRequest](c)
	if err != nil {
		return err
	}

	userId := c.Locals("user_id").(uint64)

	resp, err := h.service.UpdateLook(body, uint(lookId), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetCuratorAccountDetail Get curator account detail
//
//	@Summary		Get curator account detail
//	@Description	This endpoint is used for get curator account detail
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			curator_id	path		string	true	"Curator ID"
//	@Success		200			{object}	dto.AccountDetailResponse
//	@Failure		400			{object}	fiber.Error
//	@Router			/curator/{curator_id}/account-detail [get]
func (h *CuratorBOHandlers) GetCuratorAccountDetail(c *fiber.Ctx) error {
	curatorIDStr := c.Params("curator_id")

	curatorID, err := strconv.ParseUint(curatorIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	curators, err := h.service.GetCuratorAccountDetail(uint(curatorID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(curators)
}

// UpdateCuratorAccountDetail Update curator account detail
//
//	@Summary		Update curator account detail
//	@Description	This endpoint is used for updating curator account detail
//	@Tags			CuratorBO
//	@Security		BearerAuth
//	@Produce		application/json
//	@Param			curator_id	path		string							true	"Curator ID"
//	@Param			account		body		dto.UpdateAccountDetailRequest	true	"Account Detail"
//	@Success		200			{object}	dto.AccountDetailResponse
//	@Failure		400			{object}	fiber.Error
//	@Failure		500			{object}	fiber.Error
//	@Router			/curator/{curator_id}/account-detail [put]
func (h *CuratorBOHandlers) UpdateCuratorAccountDetail(c *fiber.Ctx) error {
	curatorIDStr := c.Params("curator_id")

	curatorID, err := strconv.ParseUint(curatorIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid curator ID"})
	}

	var request dto.UpdateAccountDetailRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	updatedAccount, err := h.service.UpdateCuratorAccountDetail(uint(curatorID), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updatedAccount)
}
