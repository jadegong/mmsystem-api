package middleware

import (
	"strings"

	"github.com/labstack/echo"
)

var (
	NoneAuthUrls map[string][]string = map[string][]string{
		"ALL": []string{"/login", "/register", "/forget", "/reset", "/captcha"},
		"GET": []string{"/file/view"}, //TODO change to the api of material detection result document
	}
)

// TODO add middleware of token
// Generate token and Store token while login(alg?)
// Another request:
// 1. Validate token from header(Get token, read from storage -- redis or memory, validate token); alg to decode token???
// 2. While valid, store Context variable, such as c.store["user"] = $UserID, get by c.Get(key)
// 3. Done
// Methods: GenerateToken, DecodeTokenInfo, SaveToken, ReadToken, GetTokenFromHeader, ValidateToken??,
// SaveVariable(could be done in middleware func)


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
