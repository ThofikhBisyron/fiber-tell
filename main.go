package main

import (
	"log"
	"os"
	"tell-be/lib"
	"tell-be/routers"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}

	// Connect database
	db, err := lib.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Create Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
		},
		AllowCredentials: true,
	}))

	routers.RouterCombine(app, db)

	// Test route
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Go Tell API",
		})
	})

	// Health check
	app.Get("/health", func(c fiber.Ctx) error {
		if err := db.Ping(c.Context()); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":   "error",
				"database": "disconnected",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "connected",
		})
	})

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}

}
