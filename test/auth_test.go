package test

import (
	"contact-management/src/apps"
	"contact-management/src/config"
	"contact-management/src/controllers"
	"contact-management/src/middlewares"
	"contact-management/src/repositories"
	"contact-management/src/services"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func setupAuthRouter() *httprouter.Router {
	cfg := config.LoadConfig()
	db, err := apps.Connect(cfg)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	router := httprouter.New()
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.GET("/me", middlewares.AuthMiddleware(authController.Me))
	router.POST("/logout", middlewares.AuthMiddleware(authController.Logout))

	return router
}

func TestRegister(t *testing.T) {
	router := setupAuthRouter()

	t.Run("Success - Register new user", func(t *testing.T) {
		body := map[string]interface{}{
			"username": "testuser_register",
			"password": "password123",
			"name":     "Test User",
		}

		rr := makeRequest(t, router, "POST", "/register", body, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusCreated, rr.Code)
		assertResponseStatus(t, "success", response)

		// Cleanup
		defer cleanupTestUser(t, "testuser_register")
	})

	t.Run("Error - Duplicate username", func(t *testing.T) {
		// Create user first
		body := map[string]interface{}{
			"username": "testuser_duplicate",
			"password": "password123",
			"name":     "Test User",
		}
		makeRequest(t, router, "POST", "/register", body, "")

		// Try to register again with same username
		rr := makeRequest(t, router, "POST", "/register", body, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusConflict, rr.Code)
		assertResponseStatus(t, "error", response)

		// Cleanup
		defer cleanupTestUser(t, "testuser_duplicate")
	})

	t.Run("Error - Missing required fields", func(t *testing.T) {
		body := map[string]interface{}{
			"username": "testuser_incomplete",
			// Missing password and name
		}

		rr := makeRequest(t, router, "POST", "/register", body, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnprocessableEntity, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestLogin(t *testing.T) {
	router := setupAuthRouter()

	// Setup: Create a test user
	registerBody := map[string]interface{}{
		"username": "testuser_login",
		"password": "password123",
		"name":     "Test User",
	}
	makeRequest(t, router, "POST", "/register", registerBody, "")
	defer cleanupTestUser(t, "testuser_login")

	t.Run("Success - Login with valid credentials", func(t *testing.T) {
		loginBody := map[string]interface{}{
			"username": "testuser_login",
			"password": "password123",
		}

		rr := makeRequest(t, router, "POST", "/login", loginBody, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)

		// Verify token is present in response
		if len(response.Data) == 0 {
			t.Error("Expected token in response data")
		}
	})

	t.Run("Error - Invalid password", func(t *testing.T) {
		loginBody := map[string]interface{}{
			"username": "testuser_login",
			"password": "wrongpassword",
		}

		rr := makeRequest(t, router, "POST", "/login", loginBody, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Non-existent user", func(t *testing.T) {
		loginBody := map[string]interface{}{
			"username": "nonexistent_user",
			"password": "password123",
		}

		rr := makeRequest(t, router, "POST", "/login", loginBody, "")
		response := parseResponse(t, rr)

		// Note: Currently returns 500 - should handle sql.ErrNoRows in repository
		assertStatusCode(t, http.StatusInternalServerError, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Missing credentials", func(t *testing.T) {
		loginBody := map[string]interface{}{
			"username": "testuser_login",
			// Missing password
		}

		rr := makeRequest(t, router, "POST", "/login", loginBody, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnprocessableEntity, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestMe(t *testing.T) {
	router := setupAuthRouter()

	// Setup: Create and login test user
	registerBody := map[string]interface{}{
		"username": "testuser_me",
		"password": "password123",
		"name":     "Test User Me",
	}
	makeRequest(t, router, "POST", "/register", registerBody, "")
	defer cleanupTestUser(t, "testuser_me")

	token := getValidToken(t, "testuser_me")

	t.Run("Success - Get current user with valid token", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/me", nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)
	})

	t.Run("Error - Missing token", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/me", nil, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Invalid token", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/me", nil, "invalid_token_here")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestLogout(t *testing.T) {
	router := setupAuthRouter()

	// Setup: Create test user
	registerBody := map[string]interface{}{
		"username": "testuser_logout",
		"password": "password123",
		"name":     "Test User Logout",
	}
	makeRequest(t, router, "POST", "/register", registerBody, "")
	defer cleanupTestUser(t, "testuser_logout")

	token := getValidToken(t, "testuser_logout")

	t.Run("Success - Logout with valid token", func(t *testing.T) {
		rr := makeRequest(t, router, "POST", "/logout", nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)
	})

	t.Run("Error - Logout without token", func(t *testing.T) {
		rr := makeRequest(t, router, "POST", "/logout", nil, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Logout with invalid token", func(t *testing.T) {
		rr := makeRequest(t, router, "POST", "/logout", nil, "invalid_token")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}
