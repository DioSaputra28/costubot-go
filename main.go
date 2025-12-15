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
	_ = services.NewCategoryService(categoryRepo)
	// _ := controllers.

	port := ":8000"
	logger.Info("Server running on port " + port)
	logger.Fatal(http.ListenAndServe(port, router))
}
