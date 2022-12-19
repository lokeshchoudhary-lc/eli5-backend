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
	app.Use(cors.New(cors.Config{AllowCredentials: true, AllowOrigins: "http://127.0.0.1:5173,http://localhost:5173,http://localhost,http://127.0.0.1,http://eli5.club,https://eli5.club"}))
}

func Create() *fiber.App {
	database.Connect()
	// seed.SeedDatabase()

	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	app := fiber.New(fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	})

	// app := fiber.New(fiber.Config{
	// 	// Override default error handler
	// 	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
	// 		if e, ok := err.(*utils.Error); ok {
	// 			return ctx.Status(e.Status).JSON(e)
	// 		} else if e, ok := err.(*fiber.Error); ok {
	// 			return ctx.Status(e.Code).JSON(utils.Error{Status: e.Code, Code: "internal-server", Message: e.Message})
	// 		} else {
	// 			return ctx.Status(500).JSON(utils.Error{Status: 500, Code: "internal-server", Message: err.Error()})
	// 		}
	// 	},
	// })

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
		if err := app.Shutdown(); err != nil {
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
