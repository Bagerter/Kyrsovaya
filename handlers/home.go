package handlers

import (
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}
	username := session.Get("username")
	if username == nil {
		return c.Render("index", fiber.Map{})
	} else {
		return c.Redirect("/dash")
	}
}
