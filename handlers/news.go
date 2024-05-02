package handlers

import (
	"context"
	"newproject/db"

	"github.com/gofiber/fiber/v2"
)

func DeleteNewsItem(c *fiber.Ctx) error {
	// Проверка, является ли пользователь администратором
	session, err := db.Sessions.Get(c)
	if err != nil || session.Get("username") != "artem" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Получение ID новости из параметра URL
	newsID := c.Params("id")
	if newsID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("News ID is required")
	}

	// Удаление новости из базы данных
	commandTag, err := db.DBConn.Exec(context.Background(), `DELETE FROM NewsItems WHERE id = $1`, newsID)
	if err != nil || commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete news item")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "News item deleted successfully"})
}
func AddNewsItem(c *fiber.Ctx) error {
	// Проверка, является ли пользователь администратором
	session, err := db.Sessions.Get(c)
	if err != nil || session.Get("username") != "artem" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	// Парсинг тела запроса
	type RequestBody struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot parse JSON")
	}

	// Вставка новости в базу данных
	commandTag, err := db.DBConn.Exec(context.Background(), `INSERT INTO NewsItems (title, content) VALUES ($1, $2)`, body.Title, body.Content)
	if err != nil || commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to add news item")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "News item added successfully"})
}
