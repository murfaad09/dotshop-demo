package routes

import (
	"github.com/gofiber/fiber/v2"

	handler "github.com/harishash/dotshop-be/integration/vndr/convictional/handler"
)

func InitConvictionalRoutes(app *fiber.App) {
	// Initialize handlers
	productHandler := handler.NewProductHandler()
	//Initialize JWT
	//	jwt := middlewares.NewAuthorizationMiddleware()

	//version the API for change tracking
	api := app.Group("/vendor")
	v1 := api.Group("/buyer")

	// Define routes
	v1.Get("/products/:id", productHandler.GetProductByID)
	v1.Get("/products", productHandler.ListProducts)
	v1.Get("/products/:product_id/images/:image_id", productHandler.DeleteProductImage)

	// Add more routes as needed
}
