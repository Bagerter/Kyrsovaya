package routes

import (
	"newproject/handlers"

	"github.com/gofiber/fiber/v2"
)

func Set(app *fiber.App) {
	request := app.Group("/req")
	app.Get("/", handlers.Home)
	app.Get("/dash", handlers.Dashboard)
	//app.Get("/collect_browser_data", handlers.ShowBrowserDataPage)
	app.Get("/captcha/verified", handlers.VerifyCaptcha)
	app.Get("/logout", handlers.Logout)
	app.Get("/domains", handlers.Domains)
	app.Get("/domain-stats/:id", handlers.DomainStats)
	request.Post("/signup", handlers.RequestSignup)
	request.Post("/login", handlers.RequestLogin)
	request.Post("/news", handlers.AddNewsItem)
	request.Post("/news", handlers.AddNewsItem)
	request.Post("/domains", handlers.AddDomain)
	request.Get("/domains", handlers.GetDomains)
	request.Delete("/domains/:id", handlers.DeleteDomain)
	request.Put("/domains/:id", handlers.UpdateDomain)
	request.Delete("/news/:id", handlers.DeleteNewsItem)
	app.Post("/collect_browser_data", handlers.CollectBrowserData)
}
