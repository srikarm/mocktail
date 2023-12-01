package main

import (
	"fmt"
	"log"
	"mocktail-api/core"
	"mocktail-api/database"
	"mocktail-api/mocktail"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Static("/", "./build")
	coreApi := app.Group("/core/v1")
	coreApi.Get("/apis", core.GetApis)
	coreApi.Get("/export", core.ExportApis)
	coreApi.Post("/api", core.CreateApi)
	coreApi.Post("/import", core.ImportApis)
	coreApi.Delete("/api/:id", core.DeleteApiByKey)

	mocktailApi := app.Group("/mocktail")
	mocktailApi.Get("/:endpoint/*", mocktail.MockApiHandler)
	mocktailApi.Post("/:endpoint/*", mocktail.MockApiHandler)
	mocktailApi.Put("/:endpoint/*", mocktail.MockApiHandler)
	mocktailApi.Patch("/:endpoint/*", mocktail.MockApiHandler)
	mocktailApi.Delete("/:endpoint/*", mocktail.MockApiHandler)

}

func initDatabase() {
	var err error
	// Read PostgreSQL connection details from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "apis")
	sslMode := getEnv("SSL_MODE", "disable")

	// Construct the connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connected to PostgreSQL database successfully!")
	
	err = database.DBConn.AutoMigrate(&core.Api{})
	if err != nil {
		log.Fatal("Error migrating database:", err)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func main() {
	svcPort := getEnv("DB_PORT", "4000") 
	app := fiber.New()
	app.Use(cors.New())

	initDatabase()
//	defer database.DBConn.Close()

	setupRoutes(app)

	log.Fatal(app.Listen(":" + svcPort))
}
