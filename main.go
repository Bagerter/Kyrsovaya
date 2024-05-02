package main

import (
	"newproject/app"
	"newproject/db"
	"newproject/routes"
)

func main() {
	db.InitPg("postgres://postgres:root@localhost:5432/protect")
	db.InitRd()
	application := app.Create()
	routes.Set(application)
	application.ListenTLS(":443", "localhost.crt", "localhost.key")
	application.Listen(":8080")
}
