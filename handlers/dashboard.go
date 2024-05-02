package handlers

import (
	"fmt"
	"newproject/db"
	"newproject/models"

	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}
	user := session.Get("username")
	if user != nil {
		userID, err := db.GetUserIDByUsername(user.(string))
		if err != nil {
			// Обработка ошибки, например, отправка ответа об ошибке
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user ID")
		}
		session.Set("user_id", userID)
		if err := session.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot save session")
		}
		userStr, ok := user.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Session username is not a string")
		}
		admin := "artem"
		isAdmin := (admin == userStr)
		newsItems, err := db.GetAllNews()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get news: %v", err))
		}
		userscount, err := db.GetUserCount()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get count users: %v", err))
		}
		totalAttacksLast7Days, err := db.GetAttackCountsForLast7Days()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get 7d attacks: %v", err))
		}

		totalAttacks, err := db.GetTotalAttackCount()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get total attacks: %v", err))
		}
		labels := db.GetLast7DaysLabels()
		fmt.Println(labels, totalAttacksLast7Days)
		dashboardData := models.DashboardData{
			IsAdmin:      isAdmin,
			Username:     userStr,
			TotalUsers:   userscount,
			TotalAttacks: totalAttacks,
			OnlineBots:   0,
			NewsItems:    newsItems,
			GraphData: models.GraphData{
				Labels: labels,
				Data:   totalAttacksLast7Days,
			},
		}

		return c.Render("1", dashboardData)
	} else {
		return c.Redirect("/")
	}
}
