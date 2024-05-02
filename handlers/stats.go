package handlers

import (
	"encoding/json"
	"fmt"
	"newproject/db"
	"newproject/models"

	"github.com/gofiber/fiber/v2"
)

func DomainStats(c *fiber.Ctx) error {
	session, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get session")
	}
	userid := session.Get("user_id")
	if userid != nil {
		iddomain, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain ID"})
		}
		gd, err := db.DomainBelongsToUser(iddomain, userid.(int))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid domain ID"})
		}
		if !gd {
			return c.Redirect("/domains")
		}
		attackinfo, err := db.GetDomainAttacks(iddomain)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot get attack info")
		}
		// domainvisits, err := db.GetDomainVisits(iddomain)
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).SendString("Cannot get domain visits")
		// }
		visitorips, err := db.GetVisitorIPs(iddomain)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot get domain visits")
		}
		attackips, err := db.GetAttackerIPs(iddomain)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot get domain visits")
		}
		country_attack, err := db.GetIPCountsByCountry(iddomain, "visitor_ips")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot get country attack visits")
		}
		country_visits, err := db.GetIPCountsByCountry(iddomain, "attacker_ips")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot get country visits")
		}
		visitCountriesDataJSON, err := json.Marshal(country_visits)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot json country visits")
		}
		attackCountriesDataJSON, err := json.Marshal(country_attack)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot json attack visits")
		}
		fmt.Println(string(visitCountriesDataJSON))
		statisticData := models.StatisticData{
			DomainName:              "Название вашего домена",
			TotalAttacks:            CountTotalAttacks(attackinfo),
			ReflectedAttacks:        CountTotalAttacks(attackinfo),
			MaxAttackPower:          FindMaxPower(attackinfo),             // Предполагается, что вы вычислили это значение
			AvgRequestsPerIP:        FindMaxRequestsPerMinute(attackinfo), // Также предполагается вычисление
			VisitCountriesDataJSON:  string(visitCountriesDataJSON),
			AttackCountriesDataJSON: string(attackCountriesDataJSON),
			VisitorIPs:              visitorips,
			AttackerIPs:             attackips,
		}
		return c.Render("stats", statisticData)
	} else {
		return c.Redirect("/")
	}
}

func FindMaxPower(attacks []models.AttackInfo) int {
	maxPower := 0
	for _, attack := range attacks {
		if attack.MaxPower > maxPower {
			maxPower = attack.MaxPower
		}
	}
	return maxPower
}

// Функция для нахождения максимального значения requests_per_minute из списка атак
func FindMaxRequestsPerMinute(attacks []models.AttackInfo) int {
	maxRequests := 0
	for _, attack := range attacks {
		if attack.RequestsPerMinute > maxRequests {
			maxRequests = attack.RequestsPerMinute
		}
	}
	return maxRequests
}

// Функция для подсчета общего количества атак
func CountTotalAttacks(attacks []models.AttackInfo) int {
	return len(attacks)
}
