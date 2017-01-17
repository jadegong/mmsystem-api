package middleware

import "github.com/labstack/echo"

var (
	NoneAuthUrls map[string][]string = map[string][]string{
		"ALL": {""},
	}
)

func TokenAuthSkipper(c echo.Context) bool {
	return true
}
