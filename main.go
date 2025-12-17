package main

import (
	"contact-management/src/apps"
	"contact-management/src/config"
	"contact-management/src/controllers"
	"contact-management/src/middlewares"
	"contact-management/src/repositories"
	"contact-management/src/services"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	cfg := config.LoadConfig()

	logger := apps.LoggingApp()
	logger.Info("Application started")

	db, err := apps.Connect(cfg)
	if err != nil {
		logger.Fatal("Database connection failed: ", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)

	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	router := httprouter.New()

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.GET("/me", middlewares.AuthMiddleware(authController.Me))
	router.POST("/logout", middlewares.AuthMiddleware(authController.Logout))

	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	router.GET("/users", middlewares.AuthMiddleware(userController.GetUser))
	router.POST("/users", middlewares.AuthMiddleware(userController.CreateUser))
	router.GET("/users/:username", middlewares.AuthMiddleware(userController.GetUserByUsername))
	router.PUT("/users/:username", middlewares.AuthMiddleware(userController.UpdateUser))
	router.DELETE("/users/:username", middlewares.AuthMiddleware(userController.DeleteUser))

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	router.GET("/categories", middlewares.AuthMiddleware(categoryController.GetAllCategories))
	router.POST("/categories", middlewares.AuthMiddleware(categoryController.CreateCategory))
	router.GET("/categories/:id", middlewares.AuthMiddleware(categoryController.GetCategoryByID))
	router.PUT("/categories/:id", middlewares.AuthMiddleware(categoryController.UpdateCategory))
	router.DELETE("/categories/:id", middlewares.AuthMiddleware(categoryController.DeleteCategory))
	
	brandProductRepo := repositories.NewBrandProductRepository(db)
	brandProductService := services.NewBrandProductService(brandProductRepo)
	brandProductController := controllers.NewBrandProductController(brandProductService)

	router.GET("/brand-products", middlewares.AuthMiddleware(brandProductController.GetAllBrandProducts))
	router.POST("/brand-products", middlewares.AuthMiddleware(brandProductController.CreateBrandProduct))
	router.GET("/brand-products/:id", middlewares.AuthMiddleware(brandProductController.GetBrandProductByID))
	router.PUT("/brand-products/:id", middlewares.AuthMiddleware(brandProductController.UpdateBrandProduct))
	router.DELETE("/brand-products/:id", middlewares.AuthMiddleware(brandProductController.DeleteBrandProduct))

	port := ":8080"
	logger.Info("Server running on port " + port)
	logger.Fatal(http.ListenAndServe(port, router))
}
