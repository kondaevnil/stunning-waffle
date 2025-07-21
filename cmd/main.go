package main

import (
	"log"
	"os"
	"vk/ecom/internal/handler"

	// "vk/ecom/internal/repository/memory"
	"vk/ecom/internal/database"
	"vk/ecom/internal/repository/postgres"
	"vk/ecom/internal/service"

	"github.com/gin-gonic/gin"
)

func setupRoutes(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/api/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}

	listings := router.Group("/api/listings")
	listings.Use(handler.OptionalAuthMiddleware())
	{
		listings.GET("/", handler.GetListings)
	}

	protected := router.Group("/api")
	protected.Use(handler.AuthMiddleware())
	{
		protected.POST("/listings", handler.CreateListing)
	}

	return router
}

func main() {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "ecom"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.MigrateDB(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userRepo := postgres.NewUserRepository(db)
	listingRepo := postgres.NewListingRepository(db)

	// userRepo := memory.NewInMemoryUserRepository()
	// listingRepo := memory.NewInMemoryListingRepository()

	authService := service.NewAuthService(userRepo)
	listingService := service.NewListingService(listingRepo, userRepo)

	// addTestData(authService, listingService)

	handler := handler.NewHandler(authService, listingService)

	router := setupRoutes(handler)

	router.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
