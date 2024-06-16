package main

import (
	"fmt"
	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nishanthrk/aspire-lms/app/configs"
	"github.com/nishanthrk/aspire-lms/app/database"
	"github.com/nishanthrk/aspire-lms/app/logger"
	"github.com/nishanthrk/aspire-lms/app/routes"
	"github.com/nishanthrk/aspire-lms/migrations"
)

func main() {
	// Load configuration from the environment file or other sources
	configs.LoadLocalConfig()
	config := configs.GetConfig()

	// Initialize the logger for structured logging
	logger.InitLogger()

	// Create a new Fiber application instance with a body size limit of 12 MB
	app := fiber.New(fiber.Config{BodyLimit: 12 * 1024 * 1024})

	// Connect to the MySQL database
	database.ConnectMysql()

	// Run database migrations
	err := migrations.RunMigrations()
	if err != nil {
		logger.Sugar.Fatal("Issue during migration 🧘🏻: ", err.Error())
	}

	// Setup middleware
	// Logger middleware for logging HTTP requests and responses
	app.Use(fiberLogger.New())
	// CORS middleware for handling Cross-Origin Resource Sharing
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, " +
			"Authorization, Access-Control-Allow-Headers, X-Platform, X-IP, X-Forwarded-For",
		AllowMethods: "*",
	}))

	// Setup New Relic monitoring if a license key is provided
	if len(configs.GetConfig().NewRelicLicense) != 0 {
		cfg := fibernewrelic.Config{
			License: configs.GetConfig().NewRelicLicense,
			AppName: configs.GetConfig().GetTenantName(),
			Enabled: true,
		}
		app.Use(fibernewrelic.New(cfg))
	}

	// Setup application routes
	routes.SetupRoutes(app)

	// Start the application and listen on the configured port
	port := fmt.Sprintf(":%s", config.Port)
	err = app.Listen(port)
	if err != nil {
		logger.Sugar.Fatal(err)
		return
	}
}
