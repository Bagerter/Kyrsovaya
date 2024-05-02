package handlers

import (
	"fmt"
	"newproject/config"
	"newproject/db"
	"newproject/firewall"

	"github.com/gofiber/fiber/v2"
)

func CollectBrowserData(c *fiber.Ctx) error {
	data := &struct {
		CookieEnabled       bool     `json:"cookieEnabled"`
		DeviceMemory        float64  `json:"deviceMemory"`
		DoNotTrack          string   `json:"doNotTrack"`
		HardwareConcurrency int      `json:"hardwareConcurrency"`
		Language            string   `json:"language"`
		Languages           []string `json:"languages"`
		MaxTouchPoints      int      `json:"maxTouchPoints"`
		Platform            string   `json:"platform"`
		UserAgent           string   `json:"userAgent"`
		Vendor              string   `json:"vendor"`
		Width               int      `json:"width"`
		Height              int      `json:"height"`
		ColorDepth          int      `json:"colorDepth"`
		PixelDepth          int      `json:"pixelDepth"`
		TimezoneOffset      int      `json:"timezoneOffset"`
		Timezone            string   `json:"timezone"`
		TouchSupport        int      `json:"touchSupport"`
		DevicePixelRatio    float64  `json:"devicePixelRatio"`
		Canvas              string   `json:"canvas"`
		WebGL               string   `json:"webgl"`
		WebglInfo           string   `json:"webglInfo"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON", "details": err.Error()})
	}
	ip, country, asn, _ := firewall.GetIpInfoforme(c)
	fraudscore := 0
	if !data.CookieEnabled {
		fraudscore = 5
	}
	if data.DeviceMemory == 0 {
		fraudscore = 5
	}
	if data.DoNotTrack == "" {
		fraudscore++
	}
	if data.HardwareConcurrency == 0 {
		fraudscore = 5
	}
	if data.Language == "" {
		fraudscore = 5
	}
	if len(data.Languages) == 0 {
		fraudscore = 5
	}
	if data.Platform == "" {
		fraudscore = 5
	}
	if data.UserAgent == "" {
		fraudscore = 5
	}
	if data.Vendor == "" {
		fraudscore = 5
	}
	if data.Width == 0 {
		fraudscore = 5
	}
	if data.Height == 0 {
		fraudscore = 5
	}
	if data.ColorDepth == 0 {
		fraudscore = 5
	}
	if data.PixelDepth == 0 {
		fraudscore = 5
	}
	if data.TimezoneOffset == 0 {
		fraudscore = 5
	}
	if data.Timezone == "" {
		fraudscore = 5
	}
	if data.DevicePixelRatio == 0 {
		fraudscore = 5
	}
	if data.Canvas == "" {
		fraudscore = 5
	}
	if data.WebGL == "" {
		fraudscore = 5
	}
	fmt.Println(data.Timezone)
	// TODO: добавить проверку на аутентичность canvas и webgl

	// TODO: сохранить fraudscore в базе данных
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
	}
	s, err := db.Sessions.Get(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
	}
	if fraudscore > 3 {
		s.Set("ban", true)
		if err := s.Save(); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
		}
		domain := c.Hostname()
		domain_info, _ := db.GetDomainFromCacheByName(domain)
		db.AddToCache(domain_info.ID, ip, asn, country, true)
		config.IncrementCounter("fail")
		return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
	} else {
		// Установить значение в сессии
		s.Set("fingerprint_collected", true)

		// Сохранить сессию
		if err := s.Save(); err != nil {
			return err
		}
		db.GetCanvas(ip, data.Canvas)
		fmt.Println(fraudscore)
		return c.JSON(fiber.Map{"status": "success"})
	}
}
