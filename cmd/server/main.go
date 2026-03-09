package main

import (
	"net/http"
	"os"
	"time"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/handler"
	"wine-shop-api/internal/middleware"
	"wine-shop-api/internal/repository"
	"wine-shop-api/internal/service"
	"wine-shop-api/pkg/config"
	"wine-shop-api/pkg/logger"
	"wine-shop-api/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "wine-shop-api/docs"
)

// @title           Wine Shop API
// @version         1.0
// @description     A RESTful API for an online wine shop.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load .env (optional — won't fail if missing in production)
	if err := godotenv.Load(); err != nil {
		// Will use structured logger once initialized
	}

	// Load and validate configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Logger not initialized yet, use stderr
		_, _ = os.Stderr.WriteString("FATAL: Configuration error: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Initialize structured logger
	logger.Init(cfg.GinMode)
	logger.Log.Info().Msg("Starting Wine Shop API")

	// Connect to Database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	logger.Log.Info().Msg("Database connected successfully")

	// Auto Migrate
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Review{},
	); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// ──────────── Repository Layer ────────────
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// ──────────── Service Layer ────────────
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartService)
	reviewService := service.NewReviewService(reviewRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)

	// ──────────── Handler Layer ────────────
	authHandler := &handler.AuthHandler{Service: userService}
	productHandler := &handler.ProductHandler{Service: productService}
	cartHandler := &handler.CartHandler{Service: cartService}
	orderHandler := &handler.OrderHandler{Service: orderService}
	reviewHandler := &handler.ReviewHandler{Service: reviewService}
	analyticsHandler := &handler.AnalyticsHandler{Service: analyticsService}

	// Initialize Cloudinary Service (optional)
	var uploadHandler *handler.UploadHandler
	cloudinaryService, err := service.NewCloudinaryService()
	if err == nil {
		uploadHandler = &handler.UploadHandler{
			CloudinaryService: cloudinaryService,
		}
		logger.Log.Info().Msg("Cloudinary service initialized")
	} else {
		logger.Log.Warn().Msg("Cloudinary not configured — image upload disabled")
	}

	// Initialize Gin engine
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "https://*.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))

	// Global Middleware
	r.Use(gin.Recovery())

	// Rate Limiters
	generalLimiter := middleware.NewRateLimiter(100, time.Minute)
	authLimiter := middleware.NewRateLimiter(10, time.Minute)

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	public := r.Group("/api")
	public.Use(middleware.RateLimitMiddleware(generalLimiter))
	{
		public.POST("/register", middleware.RateLimitMiddleware(authLimiter), authHandler.Register)
		public.POST("/login", middleware.RateLimitMiddleware(authLimiter), authHandler.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Wine Shop API is running",
			})
		})

		public.GET("/products", productHandler.GetAllProducts)
		public.GET("/products/:id", productHandler.GetProduct)
		public.GET("/products/:id/reviews", reviewHandler.GetProductReviews)
	}

	// Protected Routes (Admin)
	protectedAdmin := r.Group("/api/admin")
	protectedAdmin.Use(middleware.AdminMiddleware())
	{
		protectedAdmin.GET("/profile", func(c *gin.Context) {
			userID, _ := utils.ExtractTokenID(c)
			c.JSON(http.StatusOK, gin.H{"message": "Admin access granted", "user_id": userID})
		})

		protectedAdmin.POST("/products", productHandler.CreateProduct)
		protectedAdmin.PUT("/products/:id", productHandler.UpdateProduct)
		protectedAdmin.DELETE("/products/:id", productHandler.DeleteProduct)

		if uploadHandler != nil {
			protectedAdmin.POST("/upload", uploadHandler.UploadImage)
		}

		protectedAdmin.GET("/analytics/stats", analyticsHandler.GetDashboardStats)
		protectedAdmin.GET("/analytics/sales-by-category", analyticsHandler.GetSalesByCategory)
		protectedAdmin.GET("/analytics/top-products", analyticsHandler.GetTopProducts)
		protectedAdmin.GET("/analytics/sales-by-day", analyticsHandler.GetSalesByDay)
		protectedAdmin.GET("/analytics/recent-orders", analyticsHandler.GetRecentOrders)
	}

	// Protected Routes (User)
	protectedUser := r.Group("/api")
	protectedUser.Use(middleware.JwtAuthMiddleware())
	{
		protectedUser.GET("/me", authHandler.GetMe)
		protectedUser.POST("/cart", cartHandler.AddToCart)
		protectedUser.GET("/cart", cartHandler.GetCart)
		protectedUser.POST("/orders", orderHandler.CreateOrder)
		protectedUser.GET("/orders", orderHandler.GetOrders)
		protectedUser.POST("/products/:id/reviews", reviewHandler.CreateReview)
		protectedUser.DELETE("/products/:id/reviews/:reviewId", reviewHandler.DeleteReview)
	}

	// Start server
	logger.Log.Info().Str("port", cfg.Port).Msg("Server starting")
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server failed to start")
	}
}
