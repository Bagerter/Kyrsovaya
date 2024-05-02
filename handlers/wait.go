package handlers

import (
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func Wait() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "POST" {
			return c.Next()
		}
		// Проверка наличия сессии
		session, err := db.Sessions.Get(c)
		if err != nil {
			return err
		}
		if session.Get("fingerprint_collected") == nil {
			// Если сессии нет, то рендерим страницу сбора данных браузера
			return c.Render("browser", nil)
		}
		return c.Next()
	}
}
