package handlers

import (
	"fmt"
	"newproject/db"
	"newproject/models"

	"github.com/gofiber/fiber/v2"
)

// Add a new domain
func AddDomain(c *fiber.Ctx) error {
	var domain models.Domain
	if err := c.BodyParser(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}

	userID, ok := session.Get("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Session user ID is invalid")
	}
	domain.UserID = userID
	if err := db.AddDomain(&domain); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(domain)
}

// Update an existing domain
func DeleteDomain(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain ID"})
	}

	if err := db.DeleteDomain(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Delete an existing domain
func UpdateDomain(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		fmt.Println("domain")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain ID"})
	}

	var domain models.Domain
	if err := c.BodyParser(&domain); err != nil {
		fmt.Println("json")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}

	userID, ok := session.Get("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Session user ID is invalid")
	}
	domain.UserID = userID
	domain.ID = id // Убедитесь, что ID домена установлен правильно

	if err := db.UpdateDomain(&domain); err != nil {
		fmt.Println("db")
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(domain)
}
func GetDomains(c *fiber.Ctx) error {
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}

	userID, ok := session.Get("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Session user ID is invalid")
	}

	domains, err := db.GetAllDomains(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось получить список доменов"})
	}
	fmt.Println(domains)
	return c.JSON(domains)
}
