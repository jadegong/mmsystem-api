package middleware

import (
	"github.com/labstack/echo"
	"strings"
)

var (
	NoneAuthUrls map[string][]string = map[string][]string{
		"ALL": []string{"/login", "/register", "/forget", "/reset", "/captcha"},
		"GET": []string{"/file/view"}, //TODO change to the api of material detection result document
	}
)

// Skip request that's in NoneAuthUrls
func TokenAuthSkipper(c echo.Context) bool {
	isNoneAuthUrl := false
	if len(NoneAuthUrls["ALL"]) > 0 || len(NoneAuthUrls[c.Request().Method]) > 0 {
		urls := append(NoneAuthUrls["ALL"], NoneAuthUrls[c.Request().Method]...)
		for _, url := range urls {
			if strings.Contains(c.Request().URL.Path, url) {
				isNoneAuthUrl = true
				break
			}
		}
	}
	return isNoneAuthUrl
}
