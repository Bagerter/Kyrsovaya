package handlers

import (
	"fmt"
	"newproject/config"
	"newproject/db"
	"newproject/firewall"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func VerifyCaptcha(c *fiber.Ctx) error {
	ip, country, asn, _ := firewall.GetIpInfoforme(c)
	canvas := db.GetCanvas(ip, "empty")
	keyForCaching := ip + canvas + strconv.Itoa(time.Now().Hour())

	_, expectedSecret := db.GetCaptchaCache(keyForCaching)

	cookieName := fmt.Sprintf("%s_captcha", ip)
	userInput := c.Cookies(cookieName)

	userInputValue := strings.TrimSuffix(userInput, "{{.PublicPart}}")

	if userInputValue == expectedSecret {
		s, err := db.Sessions.Get(c)
		if err != nil {
			return c.SendString("not verified")
		}
		s.Set("captcha_passed", true)
		if err := s.Save(); err != nil {
			c.SendString("not verified")
		}
		return c.SendString("verified")
	} else {
		domain := c.Hostname()
		domain_info, _ := db.GetDomainFromCacheByName(domain)
		db.AddToCache(domain_info.ID, ip, asn, country, true)
		config.IncrementCounter("fail")
		return c.SendString("not verified")
	}
}
