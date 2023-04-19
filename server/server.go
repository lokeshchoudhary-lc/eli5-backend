package server

import (
	"eli5/api"
	"eli5/config/database"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2/middleware/recover"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(recover.New())
	// app.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: "*"}))
	// app.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: "http://localhost:5173,http://localhost"}))
	app.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: "https://eli5.club"}))
}

func Create() *fiber.App {

	// database.ConnectRedis()

	database.ConnectMongo()
	// seed.SeedDatabase()

	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	app := fiber.New(fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	})

	setupMiddlewares(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	})

	api.SetupApiRoutes(app)

	return app
}
func Listen(app *fiber.App) {

	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")

	if err := app.Listen(fmt.Sprintf("%s:%s", serverHost, serverPort)); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
