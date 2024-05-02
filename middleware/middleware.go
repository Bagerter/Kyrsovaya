package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"newproject/config"
	"newproject/db"
	"newproject/firewall"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		domain := strings.ToLower(c.Hostname())
		domain_info, _ := db.GetDomainFromCacheByName(domain)
		if domain_info == nil {
			return c.Status(fiber.StatusBadRequest).SendString("We cant find domain info")
		}
		ip, country, asn, org := firewall.GetIpInfoforme(c)
		if !domain_info.Cloudflare {
			ja3, err := firewall.GenerateJA3Hash(c)
			if err != nil || ja3 == "" {
				db.AddToCache(domain_info.ID, ip, asn, country, true)
				config.IncrementCounter("fail")
				return c.Status(fiber.StatusBadRequest).SendString("Cant generate ja3")
			}
			if !db.Fraudja3(ja3, ip) {
				db.AddToCache(domain_info.ID, ip, asn, country, true)
				config.IncrementCounter("fail")
				return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
			}
		}
		serverKeywords := []string{"cloud", "hosting", "datacenter", "proxy", "amazon", "vultr", "linode", "m247"}
		for _, keyword := range serverKeywords {
			if strings.Contains(strings.ToLower(org), keyword) {
				db.AddToCache(domain_info.ID, ip, asn, country, true)
				config.IncrementCounter("fail")
				return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
			}
		}
		excludedPaths := []string{"/captcha/verified", "/collect_browser_data"}

		for _, path := range excludedPaths {
			if strings.HasPrefix(c.Path(), path) {
				return c.Next()
			}
		}
		config.IncrementCounter("success")
		mode := config.CheckAndUpdateMode(domain_info.RateLimit, domain_info.Name)
		fmt.Println("current mode is", mode)
		s, err := db.Sessions.Get(c)
		if err != nil {
			db.AddToCache(domain_info.ID, ip, asn, country, true)
			config.IncrementCounter("fail")
			return c.Status(fiber.StatusBadRequest).SendString("Cant generate ja3")
		}
		if s.Get("ban") == true {
			db.AddToCache(domain_info.ID, ip, asn, country, true)
			config.IncrementCounter("fail")
			return c.Status(fiber.StatusBadRequest).SendString("Blocked due fraudscore")
		}
		if mode == config.ModeCaptchaBased {
			return c.Next()
		} else if mode == config.ModeCookieBased && s.Get("cookie_passed") == nil {
			s.Set("cookie_passed", true)
			if err := s.Save(); err != nil {
				db.AddToCache(domain_info.ID, ip, asn, country, true)
				config.IncrementCounter("fail")
				return c.Status(fiber.StatusBadRequest).SendString("Cant generate ja3")
			}
		} else if mode == config.ModeJSBased && s.Get("fingerprint_collected") == nil {
			return c.Render("browser", nil)
		} else if mode == config.ModeNone && s.Get("captcha_passed") == nil {
			canvas := "empty"
			canvas = db.GetCanvas(ip, canvas)
			keyForCaching := ip + canvas + strconv.Itoa(time.Now().Hour())
			publicPart := ""
			captchaData, secretPart := db.GetCaptchaCache(keyForCaching)
			if captchaData == "error" && secretPart == "error" {
				secretPart = keyForCaching[:6]
				publicPart = keyForCaching[6:]
				captchaImg := image.NewRGBA(image.Rect(0, 0, 100, 37))
				AddLabel(captchaImg, rand.Intn(90), rand.Intn(30), publicPart[:6], color.RGBA{255, 0, 0, 100})
				AddLabel(captchaImg, 25, 18, secretPart, color.RGBA{61, 140, 64, 255})

				amplitude := 2.0
				period := float64(37) / 5.0
				displacement := func(x, y int) (int, int) {
					dx := amplitude * math.Sin(float64(y)/period)
					dy := amplitude * math.Sin(float64(x)/period)
					return x + int(dx), y + int(dy)
				}
				captchaImg = WarpImg(captchaImg, displacement)

				var buf bytes.Buffer
				if err := png.Encode(&buf, captchaImg); err != nil {
					c.SendString("error code captcha")
				}
				data := buf.Bytes()

				captchaData = base64.StdEncoding.EncodeToString(data)
				error := db.CaptchaCache(keyForCaching, captchaData, secretPart)
				if error != nil {
					return c.SendString("error parsing your fingerprint")
				}

			}

			return c.Render("template", fiber.Map{
				"CaptchaData": captchaData,
				"IP":          ip,
				"PublicPart":  publicPart,
			})
		}
		if domain != "localhost" {
			originalPath := c.OriginalURL()
			targetURL := domain_info.Ip + originalPath
			return proxy.Do(c, targetURL)
		}
		return c.Next()
	}
}

func AddLabel(img *image.RGBA, x, y int, label string, color color.RGBA) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func WarpImg(src image.Image, displacement func(x, y int) (int, int)) *image.RGBA {
	bounds := src.Bounds()
	minX := bounds.Min.X
	minY := bounds.Min.Y
	maxX := bounds.Max.X
	maxY := bounds.Max.Y

	dst := image.NewRGBA(image.Rect(minX, minY, maxX, maxY))
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			dx, dy := displacement(x, y)
			if dx < minX || dx > maxX || dy < minY || dy > maxY {
				continue
			}
			dst.Set(x, y, src.At(dx, dy))
		}
	}
	return dst
}
