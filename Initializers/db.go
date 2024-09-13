package initializers

import (
	"fmt"
	"log"
	"os"

	"golang-fiber/models" // Importing the models package to use the Note struct

	"gorm.io/driver/postgres" // Importing PostgreSQL driver for GORM
	"gorm.io/gorm"            // Importing GORM for ORM functionalities
	"gorm.io/gorm/logger"     // Importing GORM logger for database query logging
)

// DB is a global variable to hold the database connection
var DB *gorm.DB

// ConnectDB initializes the database connection using the provided configuration
func ConnectDB(config *Config) {
	var err error

	// Format the Data Source Name (DSN) for PostgreSQL connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	// Open a connection to the database using the DSN
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If an error occurs, log it and exit the application
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	// Execute SQL command to create the "uuid-ossp" extension if it does not exist
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Set the GORM logger to display SQL query logs
	DB.Logger = logger.Default.LogMode(logger.Info)

	// Run database migrations for the Note model
	log.Println("Running Migrations")
	DB.AutoMigrate(&models.Note{})

	// Log successful connection
	log.Println("🚀 Connected Successfully to the Database")
}
