package main

import (
	"fmt"
	"log"
	"mocktail-api/core"
	"mocktail-api/database"
	"mocktail-api/mocktail"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/iris-contrib/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRoutes(app *iris.Application) {
	coreApi := app.Party("/core/v1")
	coreApi.Get("/apis", core.GetApis)
	coreApi.Get("/export", core.ExportApis)
	coreApi.Post("/api", core.CreateApi)
	coreApi.Post("/import", core.ImportApis)
	coreApi.Delete("/api/:id", core.DeleteApiByKey)

	mocktailApi := app.Party("/mocktail")
	mocktailApi.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Mocktail API Mocking Service"})
	})
	mocktailApi.Get("/:api", mocktail.MockApiHandler)
	mocktailApi.Post("/:api", mocktail.MockApiHandler)
	mocktailApi.Put("/:api", mocktail.MockApiHandler)
	mocktailApi.Patch("/:api", mocktail.MockApiHandler)
	mocktailApi.Delete("/:api", mocktail.MockApiHandler)

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

func notFoundHandler(ctx iris.Context) {
	ctx.StopWithJSON(iris.StatusNotFound, iris.Map{"response": "Not found"})
}

func main() {
	svcPort := getEnv("DB_PORT", "4000")
	app := iris.Default()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Accept", "content-type", "Access-Control-Allow-Origin", "X-Requested-With", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Screen", "token", "offset", "limit"},
		AllowCredentials: true,
	})
	app.UseRouter(corsOptions)

	initDatabase()
	setupRoutes(app)
	app.HandleDir("/", "./build")
	
	app.AllowMethods(iris.MethodOptions)
	// Listen for incoming HTTP/1.x & HTTP/2 clients.
	app.Run(iris.Addr(":" + svcPort))
}
