package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mendelgusmao/eulabs-api/application/rest"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*rest.JWTClaims)

		if !claims.Admin {
			return c.NoContent(http.StatusForbidden)
		}

		return next(c)
	}
}
