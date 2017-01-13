package tokenAuthMiddleware

import "github.com/labstack/echo"

func TokenAuthSkipper(c echo.Context) bool {
	return false
}
