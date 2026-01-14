package main

import (
	"log"
	"net/http"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/handler"
	"wine-shop-api/internal/middleware"
	"wine-shop-api/internal/service"
	"wine-shop-api/pkg/config"
	"wine-shop-api/pkg/utils"

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
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Initialize Gin engine
	r := gin.Default()

	// Global Middleware
	r.Use(gin.Recovery())

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

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	public := r.Group("/api")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Wine Shop API is running",
			})
		})

		// Product Routes (Public)
		public.GET("/products", productHandler.GetAllProducts)
		public.GET("/products/:id", productHandler.GetProduct)
	}

	// Protected Routes (Admin)
	protectedAdmin := r.Group("/api/admin")
	protectedAdmin.Use(middleware.JwtAuthMiddleware())
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
		// Cart Routes
		protectedUser.POST("/cart", cartHandler.AddToCart)
		protectedUser.GET("/cart", cartHandler.GetCart)

		// Order Routes
		protectedUser.POST("/orders", orderHandler.CreateOrder)
		protectedUser.GET("/orders", orderHandler.GetOrders)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
