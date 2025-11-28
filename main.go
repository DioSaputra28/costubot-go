package main

import "contact-management/src/apps"

func main() {
	logger := apps.LoggingApp()
	logger.Info("Application started")
}