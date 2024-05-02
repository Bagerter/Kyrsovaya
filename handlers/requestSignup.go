package handlers

import (
	"fmt"
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func RequestSignup(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	err := db.Signup(username, password, email)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Failed to sign up user: %v", err),
		})
	}
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Failed to sign up user: %v", err),
		})
	}
	session.Set("username", username)
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot save session")
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Sugn Up",
	})
}
