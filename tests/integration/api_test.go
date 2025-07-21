package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/handler"
	"vk/ecom/internal/repository/memory"
	"vk/ecom/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	router         *gin.Engine
	userRepo       *memory.InMemoryUserRepository
	listingRepo    *memory.InMemoryListingRepository
	authService    *service.AuthService
	listingService *service.ListingService
	handler        *handler.Handler
}

func (suite *IntegrationTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.userRepo = memory.NewInMemoryUserRepository()
	suite.listingRepo = memory.NewInMemoryListingRepository()

	suite.authService = service.NewAuthService(suite.userRepo)
	suite.listingService = service.NewListingService(suite.listingRepo, suite.userRepo)

	suite.handler = handler.NewHandler(suite.authService, suite.listingService)

	suite.router = gin.New()
	suite.setupRoutes()
}

func (suite *IntegrationTestSuite) setupRoutes() {
	auth := suite.router.Group("/api/auth")
	{
		auth.POST("/login", suite.handler.Login)
		auth.POST("/register", suite.handler.Register)
	}

	listings := suite.router.Group("/api/listings")
	listings.Use(suite.handler.OptionalAuthMiddleware())
	{
		listings.GET("/", suite.handler.GetListings)
	}

	protected := suite.router.Group("/api")
	protected.Use(suite.handler.AuthMiddleware())
	{
		protected.POST("/listings", suite.handler.CreateListing)
	}
}

func (suite *IntegrationTestSuite) TestUserRegistrationAndLogin() {
	registerReq := map[string]string{
		"login":    "testuser",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var registerResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NoError(suite.T(), err)

	user := registerResp["user"].(map[string]interface{})
	assert.Equal(suite.T(), "testuser", user["login"])
	assert.NotNil(suite.T(), user["id"])

	loginReq := map[string]string{
		"login":    "testuser",
		"password": "password123",
	}

	jsonBody, _ = json.Marshal(loginReq)
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), loginResp["token"])
	assert.Equal(suite.T(), "testuser", loginResp["login"])
	assert.NotNil(suite.T(), loginResp["user_id"])
}

func (suite *IntegrationTestSuite) TestListingFlow() {
	user := &domain.User{
		Login:    "testuser",
		Password: "password123",
	}

	registeredUser, err := suite.authService.RegisterUser(user.Login, user.Password)
	assert.NoError(suite.T(), err)

	token, _, err := suite.authService.LoginUser(user.Login, user.Password)
	assert.NoError(suite.T(), err)

	listingReq := map[string]interface{}{
		"title":       "Test Product",
		"description": "This is a test product description with enough characters",
		"image_url":   "http://example.com/image.jpg",
		"price":       100,
	}

	jsonBody, _ := json.Marshal(listingReq)
	req, _ := http.NewRequest("POST", "/api/listings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token) // Remove "Bearer " prefix as the middleware expects just the token
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NoError(suite.T(), err)

	listing := createResp["listing"].(map[string]interface{})
	assert.Equal(suite.T(), "Test Product", listing["title"])
	assert.Equal(suite.T(), registeredUser.Login, listing["author_login"])
	assert.NotNil(suite.T(), listing["id"])

	req, _ = http.NewRequest("GET", "/api/listings/", nil)
	w = httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var getResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.NoError(suite.T(), err)

	listings := getResp["listings"].([]interface{})
	assert.Len(suite.T(), listings, 1)

	firstListing := listings[0].(map[string]interface{})
	assert.Equal(suite.T(), "Test Product", firstListing["title"])
	assert.Equal(suite.T(), registeredUser.Login, firstListing["author_login"])
}

func (suite *IntegrationTestSuite) TestUnauthorizedAccess() {
	listingReq := map[string]interface{}{
		"title":       "Test Product",
		"description": "This is a test product description",
		"price":       100,
	}

	jsonBody, _ := json.Marshal(listingReq)
	req, _ := http.NewRequest("POST", "/api/listings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *IntegrationTestSuite) TestInvalidCredentials() {
	loginReq := map[string]string{
		"login":    "nonexistent",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
