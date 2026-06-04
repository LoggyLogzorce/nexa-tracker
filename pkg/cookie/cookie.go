package cookie

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Set устанавливает refresh token cookie с правильными security настройками
func Set(c *gin.Context, name, value, domain string, sameSite http.SameSite, maxAge int) {
	// Определить Secure flag в зависимости от окружения
	secure := true
	if os.Getenv("ENV") == "development" {
		secure = false
		// domain = ""
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		name,
		value,
		maxAge,
		"/",
		domain,
		secure,
		true,
	)
}

// Delete удаляет cookie
func Delete(c *gin.Context, name, domain string, sameSite http.SameSite) {
	secure := true
	if os.Getenv("ENV") == "development" {
		secure = false
		// domain = ""
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		name,
		"",
		-1, // maxAge = -1 удаляет cookie
		"/",
		domain,
		secure,
		true,
	)
}
