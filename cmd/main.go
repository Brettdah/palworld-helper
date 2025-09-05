package main

import (
	"log"

	"palworld-helper/internal/adapters/database"
	"palworld-helper/internal/core/services"
	"palworld-helper/web"
)

func main() {
	// Initialize database
	db, err := database.NewSQLiteDB("./data/palworld.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize services
	craftingService := services.NewCraftingService(db)
	adminService := services.NewAdminService(db)

	// Initialize web server
	server := web.NewServer(craftingService, adminService)

	log.Println("Palworld Helper starting on http://localhost:8080")
	log.Println("Admin interface available at http://localhost:8080/admin")

	if err := server.Start(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
