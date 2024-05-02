package handlers

import (
	"fmt"
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func RequestLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Вызываем функцию Login для проверки учетных данных пользователя
	ok, err := db.Login(username, password)
	if err != nil {
		// В случае ошибки возвращаем клиенту сообщение об ошибке
		// Можете настроить статус-код ответа в зависимости от типа ошибки
		if err.Error() == "user does not exist" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User does not exist",
			})
		} else if err.Error() == "incorrect password" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Incorrect password",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "An error occurred",
			})
		}
	}

	if !ok {
		// Если учетные данные неверны, отправляем сообщение об ошибке
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authentication failed",
		})
	}

	// Пользователь успешно аутентифицирован, отправляем подтверждение
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Не удалось получить сессию: %v", err),
		})
	}
	session.Set("username", username)
	if err := session.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot save session")
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged",
	})

}
