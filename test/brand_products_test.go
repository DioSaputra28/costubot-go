package test

import (
	"contact-management/src/apps"
	"contact-management/src/config"
	"contact-management/src/controllers"
	"contact-management/src/middlewares"
	"contact-management/src/repositories"
	"contact-management/src/services"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func setupBrandProductRouter() *httprouter.Router {
	cfg := config.LoadConfig()
	db, err := apps.Connect(cfg)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	brandProductRepo := repositories.NewBrandProductRepository(db)
	brandProductService := services.NewBrandProductService(brandProductRepo)
	brandProductController := controllers.NewBrandProductController(brandProductService)

	router := httprouter.New()
	router.GET("/brand-products", middlewares.AuthMiddleware(brandProductController.GetAllBrandProducts))
	router.POST("/brand-products", middlewares.AuthMiddleware(brandProductController.CreateBrandProduct))
	router.GET("/brand-products/:id", middlewares.AuthMiddleware(brandProductController.GetBrandProductByID))
	router.PUT("/brand-products/:id", middlewares.AuthMiddleware(brandProductController.UpdateBrandProduct))
	router.DELETE("/brand-products/:id", middlewares.AuthMiddleware(brandProductController.DeleteBrandProduct))

	return router
}

// Helper to create a test category for brand products
func createTestCategoryForBrand(t *testing.T, token string) int {
	t.Helper()

	router := setupCategoryRouter()
	body := map[string]interface{}{
		"name": "Test Category For Brands",
	}

	rr := makeRequest(t, router, "POST", "/categories", body, token)
	response := parseResponse(t, rr)

	var data map[string]interface{}
	json.Unmarshal(response.Data, &data)
	return int(data["category_id"].(float64))
}

func TestGetAllBrandProducts(t *testing.T) {
	router := setupBrandProductRouter()
	token := getValidToken(t, "testuser_brand")
	defer cleanupTestUser(t, "testuser_brand")

	t.Run("Success - Get all brand products", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/brand-products", nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)
	})

	t.Run("Error - Unauthorized without token", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/brand-products", nil, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestCreateBrandProduct(t *testing.T) {
	router := setupBrandProductRouter()
	token := getValidToken(t, "testuser_create_brand")
	defer cleanupTestUser(t, "testuser_create_brand")

	// Create test category
	categoryID := createTestCategoryForBrand(t, token)
	defer cleanupTestCategory(t, categoryID)

	t.Run("Success - Create new brand product", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Test Brand Product",
			"category_id": categoryID,
		}

		rr := makeRequest(t, router, "POST", "/brand-products", body, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusCreated, rr.Code)
		assertResponseStatus(t, "success", response)

		// Cleanup brand product
		var data map[string]interface{}
		json.Unmarshal(response.Data, &data)
		if brandProductID, ok := data["brand_product_id"].(float64); ok {
			defer cleanupTestBrandProduct(t, int(brandProductID))
		}
	})

	t.Run("Error - Invalid category_id", func(t *testing.T) {
		body := map[string]interface{}{
			"name":        "Test Brand Invalid Category",
			"category_id": 99999,
		}

		rr := makeRequest(t, router, "POST", "/brand-products", body, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusInternalServerError, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Missing required fields", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "Test Brand Missing Category",
		}

		rr := makeRequest(t, router, "POST", "/brand-products", body, token)
		response := parseResponse(t, rr)

		// Note: Currently accepts missing category_id - should add validation tag in model
		assertStatusCode(t, http.StatusInternalServerError, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestGetBrandProductByID(t *testing.T) {
	router := setupBrandProductRouter()
	token := getValidToken(t, "testuser_get_brand")
	defer cleanupTestUser(t, "testuser_get_brand")

	// Create test category and brand product
	categoryID := createTestCategoryForBrand(t, token)
	defer cleanupTestCategory(t, categoryID)

	createBody := map[string]interface{}{
		"name":        "Test Brand GetByID",
		"category_id": categoryID,
	}
	createRR := makeRequest(t, router, "POST", "/brand-products", createBody, token)
	createResponse := parseResponse(t, createRR)

	var data map[string]interface{}
	json.Unmarshal(createResponse.Data, &data)
	brandProductID := int(data["brand_product_id"].(float64))
	defer cleanupTestBrandProduct(t, brandProductID)

	t.Run("Success - Get brand product by ID", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", fmt.Sprintf("/brand-products/%d", brandProductID), nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)
	})

	t.Run("Error - Brand product not found", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/brand-products/99999", nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusNotFound, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Invalid ID format", func(t *testing.T) {
		rr := makeRequest(t, router, "GET", "/brand-products/invalid", nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusBadRequest, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestUpdateBrandProduct(t *testing.T) {
	router := setupBrandProductRouter()
	token := getValidToken(t, "testuser_update_brand")
	defer cleanupTestUser(t, "testuser_update_brand")

	// Create test category and brand product
	categoryID := createTestCategoryForBrand(t, token)
	defer cleanupTestCategory(t, categoryID)

	createBody := map[string]interface{}{
		"name":        "Test Brand Update Original",
		"category_id": categoryID,
	}
	createRR := makeRequest(t, router, "POST", "/brand-products", createBody, token)
	createResponse := parseResponse(t, createRR)

	var data map[string]interface{}
	json.Unmarshal(createResponse.Data, &data)
	brandProductID := int(data["brand_product_id"].(float64))
	defer cleanupTestBrandProduct(t, brandProductID)

	t.Run("Success - Update brand product", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"name":        "Test Brand Updated",
			"category_id": categoryID,
		}

		rr := makeRequest(t, router, "PUT", fmt.Sprintf("/brand-products/%d", brandProductID), updateBody, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)
	})

	t.Run("Error - Update non-existent brand product", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"name":        "Updated Name",
			"category_id": categoryID,
		}

		rr := makeRequest(t, router, "PUT", "/brand-products/99999", updateBody, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusNotFound, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Invalid category_id", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"name":        "Updated Name",
			"category_id": 99999,
		}

		rr := makeRequest(t, router, "PUT", fmt.Sprintf("/brand-products/%d", brandProductID), updateBody, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusInternalServerError, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}

func TestDeleteBrandProduct(t *testing.T) {
	router := setupBrandProductRouter()
	token := getValidToken(t, "testuser_delete_brand")
	defer cleanupTestUser(t, "testuser_delete_brand")

	// Create test category
	categoryID := createTestCategoryForBrand(t, token)
	defer cleanupTestCategory(t, categoryID)

	t.Run("Success - Delete brand product", func(t *testing.T) {
		// Create brand product to delete
		createBody := map[string]interface{}{
			"name":        "Test Brand To Delete",
			"category_id": categoryID,
		}
		createRR := makeRequest(t, router, "POST", "/brand-products", createBody, token)
		createResponse := parseResponse(t, createRR)

		var data map[string]interface{}
		json.Unmarshal(createResponse.Data, &data)
		brandProductID := int(data["brand_product_id"].(float64))
		defer cleanupTestBrandProduct(t, brandProductID)

		// Delete the brand product
		rr := makeRequest(t, router, "DELETE", fmt.Sprintf("/brand-products/%d", brandProductID), nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusOK, rr.Code)
		assertResponseStatus(t, "success", response)
	})

	t.Run("Error - Delete non-existent brand product", func(t *testing.T) {
		rr := makeRequest(t, router, "DELETE", "/brand-products/99999", nil, token)
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusNotFound, rr.Code)
		assertResponseStatus(t, "error", response)
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		rr := makeRequest(t, router, "DELETE", "/brand-products/1", nil, "")
		response := parseResponse(t, rr)

		assertStatusCode(t, http.StatusUnauthorized, rr.Code)
		assertResponseStatus(t, "error", response)
	})
}
