package handler

import "github.com/labstack/echo/v4"

func isAdmin(c echo.Context) bool {
	role, _ := c.Get("role").(string)
	return role == "admin"
}
