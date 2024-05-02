package handlers

import (
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func Domains(c *fiber.Ctx) error {
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}
	userid := session.Get("username")
	if userid != nil {
		return c.Render("domain", fiber.Map{})
	} else {
		return c.Redirect("/")
	}
}
