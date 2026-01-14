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
)

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
	protected := r.Group("/api/admin")
	protected.Use(middleware.JwtAuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := utils.ExtractTokenID(c)
			c.JSON(http.StatusOK, gin.H{"message": "Admin access granted", "user_id": userID})
		})

		// Product Routes (Admin)
		protected.POST("/products", productHandler.CreateProduct)
		protected.PUT("/products/:id", productHandler.UpdateProduct)
		protected.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
