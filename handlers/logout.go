package handlers

import (
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Не удалось получить сессию")
	}

	err = session.Destroy()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Не удалось удалить сессию")
	}

	return c.Redirect("/")
}
