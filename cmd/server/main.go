package main

import (
	"log"
	"net/http"
	"time"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/handler"
	"wine-shop-api/internal/middleware"
	"wine-shop-api/internal/service"
	"wine-shop-api/pkg/config"
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
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to Database
	config.ConnectDatabase()

	// Auto Migrate
	err := config.DB.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Review{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
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
			// Allow any Vercel preview or production domain
			return true
		},
	}))

	// Global Middleware
	r.Use(gin.Recovery())

	// Rate Limiter: 100 requests per minute for general routes
	generalLimiter := middleware.NewRateLimiter(100, time.Minute)
	// Rate Limiter: 10 requests per minute for auth routes (prevent brute force)
	authLimiter := middleware.NewRateLimiter(10, time.Minute)

	// Initialize Handlers
	authHandler := &handler.AuthHandler{
		Service: &service.UserService{},
	}
	productHandler := &handler.ProductHandler{
		Service: &service.ProductService{},
	}
	cartService := &service.CartService{}
	cartHandler := &handler.CartHandler{
		Service: cartService,
	}
	orderHandler := &handler.OrderHandler{
		Service: &service.OrderService{
			CartService: cartService,
		},
	}
	reviewHandler := &handler.ReviewHandler{
		Service: &service.ReviewService{},
	}

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	public := r.Group("/api")
	public.Use(middleware.RateLimitMiddleware(generalLimiter))
	{
		// Auth Routes with stricter rate limit
		public.POST("/register", middleware.RateLimitMiddleware(authLimiter), authHandler.Register)
		public.POST("/login", middleware.RateLimitMiddleware(authLimiter), authHandler.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Wine Shop API is running",
			})
		})

		// Product Routes (Public)
		public.GET("/products", productHandler.GetAllProducts)
		public.GET("/products/:id", productHandler.GetProduct)

		// Review Routes (Public - Read)
		public.GET("/products/:id/reviews", reviewHandler.GetProductReviews)
	}

	// Protected Routes (Admin) - Requires admin role
	protectedAdmin := r.Group("/api/admin")
	protectedAdmin.Use(middleware.AdminMiddleware())
	{
		protectedAdmin.GET("/profile", func(c *gin.Context) {
			userID, _ := utils.ExtractTokenID(c)
			c.JSON(http.StatusOK, gin.H{"message": "Admin access granted", "user_id": userID})
		})

		// Product Routes (Admin)
		protectedAdmin.POST("/products", productHandler.CreateProduct)
		protectedAdmin.PUT("/products/:id", productHandler.UpdateProduct)
		protectedAdmin.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Protected Routes (User)
	protectedUser := r.Group("/api")
	protectedUser.Use(middleware.JwtAuthMiddleware())
	{
		// User Info Route
		protectedUser.GET("/me", authHandler.GetMe)

		// Cart Routes
		protectedUser.POST("/cart", cartHandler.AddToCart)
		protectedUser.GET("/cart", cartHandler.GetCart)

		// Order Routes
		protectedUser.POST("/orders", orderHandler.CreateOrder)
		protectedUser.GET("/orders", orderHandler.GetOrders)

		// Review Routes (Protected - Write)
		protectedUser.POST("/products/:id/reviews", reviewHandler.CreateReview)
		protectedUser.DELETE("/products/:id/reviews/:reviewId", reviewHandler.DeleteReview)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
