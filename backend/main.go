package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/username/backend/config"
	"github.com/username/backend/controllers"
	"github.com/username/backend/models"
	"github.com/username/backend/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	// Initialize models
	userModel := models.NewUserModel(db)
	departmentModel := models.NewDepartmentModel(db)
	employeeModel := models.NewEmployeeModel(db)
	buyerModel := models.NewBuyerModel(db)

	// Initialize controllers
	authController := controllers.NewAuthController(userModel)
	userController := controllers.NewUserController(userModel)
	departmentController := controllers.NewDepartmentController(departmentModel)
	employeeController := controllers.NewEmployeeController(employeeModel)
	buyerController := controllers.NewBuyerController(buyerModel)

	// Initialize router
	router := routes.NewRouter(authController, userController, departmentController, employeeController, buyerController)
	mux := router.SetupRoutes()

	// Get server port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	// Start the server
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Database connected successfully\n")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
