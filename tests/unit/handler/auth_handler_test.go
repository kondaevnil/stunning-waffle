package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/handler"
	"vk/ecom/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestHandler_Register(t *testing.T) {
	t.Run("should successfully register user", func(t *testing.T) {
		mockAuthService := new(mocks.MockAuthService)
		mockListingService := new(mocks.MockListingService)
		h := handler.NewHandler(mockAuthService, mockListingService)

		user := &domain.User{
			ID:    1,
			Login: "testuser",
		}

		mockAuthService.On("RegisterUser", "testuser", "password123").Return(user, nil)

		router := setupTestRouter()
		router.POST("/register", h.Register)

		reqBody := map[string]string{
			"login":    "testuser",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		userResponse := response["user"].(map[string]interface{})
		assert.Equal(t, float64(1), userResponse["id"])
		assert.Equal(t, "testuser", userResponse["login"])

		mockAuthService.AssertExpectations(t)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		mockAuthService := new(mocks.MockAuthService)
		mockListingService := new(mocks.MockListingService)
		h := handler.NewHandler(mockAuthService, mockListingService)

		router := setupTestRouter()
		router.POST("/register", h.Register)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid input", response["error"])
	})

	t.Run("should fail when service returns error", func(t *testing.T) {
		mockAuthService := new(mocks.MockAuthService)
		mockListingService := new(mocks.MockListingService)
		h := handler.NewHandler(mockAuthService, mockListingService)

		mockAuthService.On("RegisterUser", "testuser", "password123").Return(nil, errors.New("user already exists"))

		router := setupTestRouter()
		router.POST("/register", h.Register)

		reqBody := map[string]string{
			"login":    "testuser",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user already exists", response["error"])

		mockAuthService.AssertExpectations(t)
	})
}

func TestHandler_Login(t *testing.T) {
	t.Run("should successfully login user", func(t *testing.T) {
		mockAuthService := new(mocks.MockAuthService)
		mockListingService := new(mocks.MockListingService)
		h := handler.NewHandler(mockAuthService, mockListingService)

		user := &domain.User{
			ID:    1,
			Login: "testuser",
		}
		token := "jwt.token.here"

		mockAuthService.On("LoginUser", "testuser", "password123").Return(token, user, nil)

		router := setupTestRouter()
		router.POST("/login", h.Login)

		reqBody := map[string]string{
			"login":    "testuser",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, token, response["token"])
		assert.Equal(t, float64(1), response["user_id"])
		assert.Equal(t, "testuser", response["login"])

		mockAuthService.AssertExpectations(t)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		mockAuthService := new(mocks.MockAuthService)
		mockListingService := new(mocks.MockListingService)
		h := handler.NewHandler(mockAuthService, mockListingService)

		router := setupTestRouter()
		router.POST("/login", h.Login)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid input", response["error"])
	})

	t.Run("should fail with invalid credentials", func(t *testing.T) {
		mockAuthService := new(mocks.MockAuthService)
		mockListingService := new(mocks.MockListingService)
		h := handler.NewHandler(mockAuthService, mockListingService)

		mockAuthService.On("LoginUser", "testuser", "wrongpassword").Return("", nil, errors.New("invalid credentials"))

		router := setupTestRouter()
		router.POST("/login", h.Login)

		reqBody := map[string]string{
			"login":    "testuser",
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid login or password", response["error"])

		mockAuthService.AssertExpectations(t)
	})
}
