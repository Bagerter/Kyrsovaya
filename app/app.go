package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func Create() *fiber.App {
	app := fiber.New(fiber.Config{
		Views: html.New("./templates", ".html"),
	})
	app.Static("/static", "./static")
	app.Use(logger.New())
	return app
}
