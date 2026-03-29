package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	"github.com/harishash/dotshop-be/integration/aws"
	"github.com/harishash/dotshop-be/integration/klaviyo"
	handlers "github.com/harishash/dotshop-be/internal/handlers"

	repository "github.com/harishash/dotshop-be/internal/repositories"
	services "github.com/harishash/dotshop-be/internal/services"
	"github.com/harishash/dotshop-be/internal/utils/auth"
	logger "github.com/harishash/dotshop-be/internal/utils/logger"

	_ "github.com/harishash/dotshop-be/internal/handlers/docs"

	config "github.com/harishash/dotshop-be/internal/config"
)

// @title						DotShop API
// @version					1.0
// @description				This is a swagger for DotShop
// @host						localhost:8080
// @BasePath					/api/v1
// @schemes					http
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func InitRoutes(app *fiber.App) {
	// get config
	cfg := config.GetConfig()

	app.Get("/swagger/*", swagger.HandlerDefault)
	// productHandler := handlers.NewProductsHandler()
	paymentHandler := handlers.NewPaymentHandler()

	// Initialize repository
	db := repository.GetDatabaseConnection().Connection

	// Initialize service
	klaviyo := klaviyo.NewKlaviyoAPI(cfg.KlaviyoKey)
	awsService, err := aws.InitializeSession(cfg.AWSREGION, cfg.AWSAccessID, cfg.AWSSecretAccessKey)
	if err != nil {
		log.Warn("Error initializing AWS session", err)
	}

	userRepo := repository.NewUserRepository(db, &klaviyo)
	orderRepo := repository.NewOrderRepository(db)

	productRepo := repository.NewProductsRepository(db)
	userService := services.NewUserService(awsService, userRepo, orderRepo, productRepo)
	userHandlers := handlers.NewUserHandler(userService)

	app.Use(cors.New())
	app.Use(logger.Fiber())
	//Initialize JWT
	jwt := auth.JWTMiddleware()

	//Permissive CORS middleware configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                      // Allow all origins
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS", // Allow all methods
		AllowHeaders: "*",                                      // Allow all headers
	}))
	// Set the db connection for handlers
	handlers.SetDependencies(db, userService)

	//version the API for change tracking
	api := app.Group("/api")
	v1 := api.Group("/v1")
	// Setup routes
	v1.Post("/auth/google/callback", handlers.GoogleSSOCallback)

	// Define routes

	v1.Post("/signup", userHandlers.CreateUser)
	v1.Post("/signin", userHandlers.Authorize)
	v1.Post("/forgot-password/send-email", userHandlers.SendForgotPasswordEmail)
	v1.Patch("/forgot-password", jwt, userHandlers.ForgotPassword)

	v1.Get("/payment/config/:gateway", paymentHandler.GetConfig)
	v1.Post("/payment/authorise/:gateway", paymentHandler.Authorize)
	v1.Post("/payment/create-payment-intent/:gateway", paymentHandler.CreatePayment)
	v1.Post("/payment/capture/:gateway", paymentHandler.CapturePayment)

	v1.Get("/validate", userHandlers.ValidateToken)

	//All routes after this needs valid token
	// app.Use(jwt)

	//User Routes
	v1.Get("/users", userHandlers.GetUsers)
	v1.Get("/users/:id",
		auth.JWTMiddleware(),
		userHandlers.GetUserByID)

	// Payment Routes
	// v1.Get("/config", jwt, paymentHandler.GetPaymentConfig)
	// v1.Post("/create-payment-intent", jwt, paymentHandler.GetPaymentConfig)
	// v1.Post("/config", jwt, paymentHandler.GetPaymentConfig)
	// v1.Post("/config", jwt, paymentHandler.GetPaymentConfig)
	// v1.Post("/config", jwt, paymentHandler.GetPaymentConfig)

	// curatorBO := app.Group("/api/v1/curatorbo/dashboard")
	// curatorStoreFront := app.Group("/api/v1/curatorstorefront")

	// Initialize repositories
	commonRepo := repository.NewCommonRepo()
	curatorBORepo := repository.NewCuratorBORepo()
	curatorOnboardingRepo := repository.NewCuratorOnboardingRepo(db, &klaviyo)
	checkoutsRepo := repository.NewCheckoutsRepo()
	cartsRepo := repository.NewCartsRepo()
	adminRepo := repository.NewAdminRepo()
	wishlistRepo := repository.NewWishlistRepository(db)
	curatorDashboardRepo := repository.NewCuratorDashboardRepository(db, productRepo)
	reviewRepo := repository.NewReviewRepository(db)
	payoutRepo := repository.NewPayoutRepository(db)
	categoryRepo := repository.NewCategoryRepository(db, productRepo)
	promotionRepo := repository.NewPromotionRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	notificationRepo := repository.NewNotificationRepository(db, adminRepo)
	adminDashboardRepo := repository.NewAdminDashboardRepository(db, productRepo)

	// Initialize services
	notificationService := services.NewNotificationService(notificationRepo)
	commonService := services.NewCommonService(commonRepo)
	curatorBOService := services.NewCuratorBOService(curatorBORepo, awsService, notificationRepo)
	curatorOnboardingService := services.NewCuratorOnboardingService(curatorOnboardingRepo)
	curatorFSService := services.NewCuratorStoreFrontService(commonRepo)

	checkoutsService := services.NewCheckoutsService(checkoutsRepo)
	cartsService := services.NewCartsService(cartsRepo)
	orderService := services.NewOrderService(orderRepo, userRepo, notificationRepo, adminRepo, awsService, productRepo)
	adminService := services.NewAdminService(adminRepo)
	reviewService := services.NewReviewService(reviewRepo)
	wishlistService := services.NewWishlistService(wishlistRepo)
	curatorDashboardService := services.NewCuratorDashboardService(curatorDashboardRepo)
	payoutService := services.NewPayoutService(payoutRepo)
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	brandService := services.NewBrandService(brandRepo, productRepo)
	promotionService := services.NewPromotionService(promotionRepo)
	adminDashboardService := services.NewAdminDashboardService(adminDashboardRepo, orderRepo, productRepo, userRepo, reviewRepo, commonRepo)
	commentService := services.NewCommentService(commentRepo)

	// Initialize handlers
	// commonHandlers := handlers.NewCommonHandlers(commonService)
	curatorBOHandlers := handlers.NewCuratorBOHandlers(curatorBOService)
	curatorStoreFrontHandlers := handlers.NewCuratorStoreFrontHandlers(commonService, curatorFSService)
	curatorOnboardingHandlers := handlers.NewCuratorOnboardingHandlers(curatorOnboardingService)
	checkoutsHandlers := handlers.NewCheckoutsHandlers(checkoutsService)
	cartsHandlers := handlers.NewCartsHandlers(cartsService)
	orderHandler := handlers.NewOrderHandler(orderService)
	adminHandler := handlers.NewAdminHandlers(adminService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	curatorDashboardHandler := handlers.NewCuratorDashboardHandler(curatorDashboardService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	payoutHandler := handlers.NewPayoutHandler(payoutService)
	productHandler := handlers.NewProductsHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	promotionHandler := handlers.NewPromotionHandler(promotionService)
	commentHandler := handlers.NewCommentHandler(commentService, productService)
	brandHandler := handlers.NewBrandHandler(brandService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	adminDashboardHandler := handlers.NewAdminDashboardHandler(adminDashboardService)

	// Common routes
	// api.Get("/collections", commonHandlers.GetAllCollections)
	api.Get("/v1/looks", curatorStoreFrontHandlers.GetAllLooks)

	// api.Get("/collections/:id/products", commonHandlers.GetAllProductsByCollectionID)
	// api.Get("/looks/:id/products", commonHandlers.GetAllProductsByLookID)

	//health
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	// Routes
	v1.Post("/notifications", jwt, notificationHandler.CreateNotification)
	v1.Get("/notifications", jwt, notificationHandler.GetNotifications)
	v1.Put("/notifications/:id/read", notificationHandler.MarkNotificationAsRead)

	v1.Post("/signup", userHandlers.CreateUser)
	v1.Post("/signin", userHandlers.Authorize)
	v1.Get("/validate", userHandlers.ValidateToken)

	v1.Get("/products", productHandler.GetAllProducts)
	v1.Get("/brands", productHandler.GetBrands)
	v1.Get("/dotshop/curator", curatorStoreFrontHandlers.GetDotShopCuratorID)

	// consumer routes
	consumer := api.Group("/v1/consumer")
	{
		consumer.Put("/profile", jwt, userHandlers.UpdateUserProfile)
		consumer.Post("/address", jwt, userHandlers.AddUserAddress)
		consumer.Get("/address", jwt, userHandlers.GetAllUserAddress)
		consumer.Get("/address/:id", jwt, userHandlers.GetUserAddressById)
		consumer.Patch("/address/:id", jwt, userHandlers.UpdateAddressById)
		consumer.Delete("/address/:id", jwt, userHandlers.DeleteAddress)
		consumer.Get("/orders", jwt, userHandlers.ConsumerOrdersList)

	}

	admin := api.Group("/v1/admin")
	{

		admin.Post("/curator/:curator_id/status", adminHandler.ChangeCuratorStatus)
		admin.Delete("/product/:ids", productHandler.DeleteProducts)

		// Category routes
		category := admin.Group("/categories")
		{
			category.Post("/", jwt, categoryHandler.AddNewCategory)
			category.Patch("/:id", jwt, categoryHandler.UpdateExistingCategory)
			category.Put("/products/:product_id", jwt, categoryHandler.UpdateProductInCategory)
			category.Put("/product/:product_id", jwt, categoryHandler.UpdateSingleProduct)
			category.Delete("/:id", jwt, categoryHandler.DeleteCategory)
			category.Get("/", categoryHandler.GetAllCategories)
			category.Get("/:id", categoryHandler.GetCategoryByID)

			subCategory := category.Group("/:category_id/subcategories")
			{
				subCategory.Post("/", jwt, categoryHandler.CreateSubCategory)
				subCategory.Post("/:id", jwt, categoryHandler.AddProductInSubCategory)
				subCategory.Patch("/:id", jwt, categoryHandler.UpdateExistingSubcategory)
				subCategory.Delete("/:id", jwt, categoryHandler.DeleteSubcategory)
			}
		}

		catalog := admin.Group("/catalog")
		{
			catalog.Get("/products", jwt, categoryHandler.GetCatalogProducts)
			catalog.Get("/brands", jwt, brandHandler.GetBrands)
			catalog.Patch("/brands/:id", jwt, brandHandler.UpdateBrandStatus)
		}

		adminDashboard := admin.Group("/dashboard")
		{
			adminDashboard.Get("/top-selling-brands", jwt, adminDashboardHandler.GetTopSellingBrands)
			adminDashboard.Get("/top-selling-products", jwt, adminDashboardHandler.GetTopSellingProducts)
			adminDashboard.Get("/top-curators", jwt, adminDashboardHandler.GetTopCurators)
			adminDashboard.Get("/sale-by-category", jwt, adminDashboardHandler.GetSaleByCategory)
			adminDashboard.Get("/top-wishlist", jwt, adminDashboardHandler.GetTopWishlistProducts)
		}

		promotions := admin.Group("/promotions")
		{
			promotions.Post("/", jwt, promotionHandler.CreatePromotion)
			promotions.Get("/", jwt, promotionHandler.GetPromotions)
			promotions.Get("/:id", jwt, promotionHandler.GetPromotionByID)
			promotions.Put("/:id", jwt, promotionHandler.UpdatePromotion)
			promotions.Delete("/:id", jwt, promotionHandler.DeletePromotion)
			promotions.Post("/discount", jwt, promotionHandler.ApplyBulkDiscount)
		}

		reviews := admin.Group("/reviews")
		{
			reviews.Get("/products", jwt, productHandler.GetAllProductsWithStats)
			reviews.Post("/comments", jwt, commentHandler.CreateComment)
			reviews.Patch("/comments/:id", jwt, commentHandler.UpdateComment)
			reviews.Delete("/comments/:id", jwt, commentHandler.DeleteComment)
			reviews.Delete("/:id", jwt, adminDashboardHandler.DeleteUserReview)
		}

		sales := admin.Group("/sales")
		{
			sales.Get("/orders", jwt, adminDashboardHandler.GetOrderSales)
			sales.Get("/returns", jwt, adminDashboardHandler.GetOrderReturns)
			sales.Patch("/return/:return_id/status", jwt, adminDashboardHandler.UpdateReturnStatus)

		}

		// Admin Customer routes
		customer := admin.Group("/customers")
		{
			customer.Get("/all", jwt, adminDashboardHandler.GetAllCustomers)
			customer.Get("/:user_id/order-list", jwt, adminDashboardHandler.ConsumerOrdersListByUserId)
			customer.Delete("/:id", jwt, adminDashboardHandler.DeleteUser)
			customer.Patch("/:user_id/block", jwt, adminDashboardHandler.BlockCustomer)

		}

		curators := admin.Group("/curators")
		{
			curators.Get("/all", jwt, adminDashboardHandler.GetAllCurators)
			curators.Get("/:curator_id/order-list", jwt, adminDashboardHandler.CuratorsOrdersListById)
			curators.Patch("/:id/block", jwt, adminDashboardHandler.BlockCurator)
			curators.Delete("/:id", jwt, adminDashboardHandler.DeleteCurator)
			curators.Get("/:id/listed-products", jwt, adminDashboardHandler.GetListedProducts)
		}

		financials := admin.Group("/financials")
		{
			financials.Get("/payment-distribution", jwt, adminDashboardHandler.GetPaymentDistribution)
		}
	}

	// Curator BO routes
	curator := api.Group("/v1/curator")
	{
		curator.Post("/addproduct", jwt, curatorBOHandlers.AddProduct)
		curator.Delete("/:curator_id/feature/product/:product_id", jwt, curatorBOHandlers.DeleteProductFromFeatureByID)
		curator.Post("/feature/:feature_id/addproduct", jwt, curatorBOHandlers.AddProductToFeature)

		curator.Post("/addcollection", jwt, curatorBOHandlers.AddCollection)
		curator.Post("/collection/:collection_id/addproduct", jwt, curatorBOHandlers.AddProductToCollection)
		curator.Put("/collection/:collection_id", jwt, curatorBOHandlers.UpdateCollection)
		curator.Delete("/collection/:collection_id", jwt, curatorBOHandlers.DeleteCollectionByID)
		curator.Delete("/collection/:collection_id/product/:product_id", jwt, curatorBOHandlers.DeleteProductFromCollectionByID)

		curator.Post("/collection/addsection", jwt, curatorBOHandlers.AddCollectionSection)
		curator.Post("/section/:section_id/addproduct", jwt, curatorBOHandlers.AddProductToSection)
		curator.Delete("/section/:section_id/product/:product_id", jwt, curatorBOHandlers.DeleteProductFromSectionByID)
		curator.Put("/section/:section_id", jwt, curatorBOHandlers.UpdateSection)

		curator.Post("/addlook", jwt, curatorBOHandlers.AddLook)
		curator.Post("/look/:look_id/addproduct", jwt, curatorBOHandlers.AddProductToLook)
		curator.Put("/look/:look_id", jwt, curatorBOHandlers.UpdateLook)

		// Search routes
		curator.Get("/look/search", curatorStoreFrontHandlers.SearchLookByName)
		curator.Get("/:curator_id/look/search/products", curatorStoreFrontHandlers.SearchProductsWithinCuratorLooks)
		curator.Get("/:curator_id/products/search", curatorStoreFrontHandlers.SearchFeatureProductsByName)
		curator.Get("/:curator_id/collection/search", curatorStoreFrontHandlers.SearchCollectionByName)
		curator.Get("/collection/:collection_id/products/search", curatorStoreFrontHandlers.SearchCollectionProductByName)

		curator.Delete("/look/:look_id", jwt, curatorBOHandlers.DeleteLookByID)
		curator.Delete("/section/:section_id", jwt, curatorBOHandlers.DeleteSectionByID)
		curator.Delete("/look/:look_id/product/:product_id", jwt, curatorBOHandlers.DeleteProductFromLookByID)

		curator.Get("/all", curatorBOHandlers.GetAllCurators)
		curator.Get("/orders", jwt, orderHandler.OrdersList)

		curator.Get("/:curator_id", curatorBOHandlers.GetCuratorWithCuratorID)

		curator.Post("/withdraw", curatorBOHandlers.Withdraw)
		// curatorBO.Get("/profile", curatorBOHandlers.GetProfile)
		curator.Put("/profile", jwt, curatorBOHandlers.UpdateProfile)
		curator.Put("/profile/password", jwt, curatorBOHandlers.ChangePassword)
		curator.Put("/:curator_id/sociallink/:link_id", jwt, curatorBOHandlers.UpdateSocialMediaLink)
		curator.Post("/:curator_id/sociallink", jwt, curatorBOHandlers.InsertSocialMediaLink)
		curator.Delete("/:curator_id/sociallink/:link_id", jwt, curatorBOHandlers.DeleteSocialMediaLink)
		curator.Get("/:curator_id/account-detail", jwt, curatorBOHandlers.GetCuratorAccountDetail)
		curator.Put("/:curator_id/account-detail", jwt, curatorBOHandlers.UpdateCuratorAccountDetail)

		// Curator Onboarding Routes
		curator.Post("/onboarding", curatorOnboardingHandlers.CreateCurator)
		curator.Get("/shopname/:shop_name", curatorOnboardingHandlers.CheckShopName)
		curator.Get("/store/:store_name", curatorOnboardingHandlers.GetCuratorByShopName)

		//  Curator Storefront routes
		// curator.Get("/:curator_id/allproducts",
		// 	jwt,
		// 	middlewares.RoleMiddleware(db, []string{constants.CURATOR}),
		// middlewares.PermissionMiddleware(db, constants.VIEW_FEATURED_PRODUCTS),
		// curatorStoreFrontHandlers.GetAllProducts)
		curator.Get("/:curator_id/allproducts", curatorStoreFrontHandlers.GetAllProducts)
		curator.Get("/:curator_id/allcollections", curatorStoreFrontHandlers.GetAllCollections)
		curator.Get("/:curator_id/alllooks", curatorStoreFrontHandlers.GetCuratorAllLooks)
		curator.Get("/collections/:collection_id/products", curatorStoreFrontHandlers.GetAllProductsByCollectionID)
		curator.Get("/section/:section_id", curatorStoreFrontHandlers.FetchSectionByID)
		curator.Get("/collection/:collection_id/section/search", curatorStoreFrontHandlers.SearchSectionByName)

		curator.Get("/looks/:look_id/products", curatorStoreFrontHandlers.GetAllProductsByLookID)

		payout := curator.Group("/payout")
		{
			payout.Get("/details", jwt, payoutHandler.GetPayoutDetails)
			payout.Get("/history", jwt, payoutHandler.GetPayoutHistory)
		}

		// Routes for dashboard pages
		dashboard := curator.Group("/:curator_id/dashboard")
		{
			dashboard.Get("/sales", jwt, curatorDashboardHandler.GetGraphDataForSales)
			dashboard.Get("/revenue", jwt, curatorDashboardHandler.GetGraphDataForRevenue)
			dashboard.Get("/orders", jwt, curatorDashboardHandler.OrderGraphData)
			dashboard.Get("/units", jwt, curatorDashboardHandler.GetGraphDataForUnits)
			dashboard.Get("/returns", jwt, curatorDashboardHandler.GraphDataForReturns)
			dashboard.Get("/average-order-value", jwt, curatorDashboardHandler.GetGraphDataForAOV)
			dashboard.Get("/average-units-per-order", jwt, curatorDashboardHandler.GetGraphDataForAUPOrder)
			dashboard.Get("/top-wishlist", jwt, curatorDashboardHandler.GetCuratorTopWishlistProducts)
			dashboard.Get("/top-selling-products", jwt, curatorDashboardHandler.GetCuratorTopSellingProducts)
			dashboard.Get("/top-selling-brands", jwt, curatorDashboardHandler.GetCuratorTopSellingBrands)
			dashboard.Get("/top-purchasers", jwt, curatorDashboardHandler.GetCuratorTopPurchasers)
			dashboard.Get("/sale-by-category", jwt, curatorDashboardHandler.GetCuratorSaleByCategory)
		}
	}

	v1.Get("/global-search", curatorStoreFrontHandlers.GlobalSearch)

	// Checkouts routes
	checkouts := api.Group("/v1/checkouts")
	{
		checkouts.Get("/product/:id", jwt, checkoutsHandlers.GetProductByID)
		// Add other routes for checkouts
	}

	//Carts routes
	carts := api.Group("/v1/cart")
	{
		//carts.Post("/buy", cartsHandlers.BuyNow)
		carts.Get("/user/:user_id", cartsHandlers.GetCartByUserID)
		carts.Post("/add", jwt, cartsHandlers.CreateCart)
		carts.Put("/:cart_id/product/:variant_id", jwt, cartsHandlers.UpdateCartItemQuantity)
		carts.Post("/:cart_id/items", jwt, cartsHandlers.AddItemsToCart)
		carts.Delete("/:cart_id", jwt, cartsHandlers.DeleteCart)
		carts.Delete("/:cart_id/items/:variant_id", jwt, cartsHandlers.DeleteCartItem)
	}
	// orders routes
	v1.Post("/order/create", jwt, orderHandler.CreateOrder)
	v1.Post("/order/return", jwt, orderHandler.CreateReturn)

	// Wishlist routes
	wishlist := api.Group("/v1/wishlist/:user_id")
	wishlist.Get("/products", jwt, wishlistHandler.GetWishlist)
	wishlist.Post("/add", jwt, wishlistHandler.AddToWishlist)
	wishlist.Delete("/remove", jwt, wishlistHandler.RemoveFromWishlist)
	// Setup routes
	users := api.Group("/v1/user")
	{
		users.Post("/reviews", jwt, reviewHandler.CreateReview)
		users.Get("/products/:product_id/reviews", reviewHandler.GetReviewsByProductID)
		users.Get("/products/:product_id/reviews/curator/:curator_id", jwt, reviewHandler.GetReviewsByProductIDAndCuratorID)
		users.Put("/reviews/:review_id", jwt, reviewHandler.UpdateReview)
		users.Delete("/reviews/:review_id", jwt, reviewHandler.DeleteReview)
		// Users routes
		users := api.Group("/v1/user")
		{
			users.Post("/order/:order_id/cancel", jwt, orderHandler.CancelOrder)
			users.Get("/profile", jwt, userHandlers.GetProfile)
		}
	}
}
