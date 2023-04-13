package main

import (
	"fmt"
	"restaurant/database"
	"restaurant/pkg/mysql"
	"restaurant/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	mysql.DatabaseInit()
	database.RunMigration()

	routes.RouteInit(e.Group("api/v1"))

	port := "5000"
	fmt.Println("server running on port", port)
	e.Logger.Fatal(e.Start("localhost:" + port))
}
