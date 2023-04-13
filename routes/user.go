package routes

import (
	"restaurant/handler"
	"restaurant/pkg/middleware"
	"restaurant/pkg/mysql"
	"restaurant/repositories"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Group) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handler.HandlerUser(userRepository)

	e.GET("/users", h.FindUsers)
	e.GET("/user/:id", h.GetUser)
	e.DELETE("/user", middleware.Auth(h.DeleteUser))
}
