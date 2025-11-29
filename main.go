package main

import (
	"contact-management/src/apps"
	"contact-management/src/config"
	"contact-management/src/repositories"
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

	userRepo = repositories.NewUserRepository(db)



}