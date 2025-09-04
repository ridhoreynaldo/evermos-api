package main

import (
    "log"
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "evermos-api/config"
    "evermos-api/routes"
)

func main() {
    // Connect to database
    config.ConnectDB()

    // Create fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "status":  code,
                "message": err.Error(),
            })
        },
    })

    // Middleware
    app.Use(cors.New())
    app.Use(logger.New())

    // Routes
    routes.SetupRoutes(app)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    log.Fatal(app.Listen(":" + port))
}