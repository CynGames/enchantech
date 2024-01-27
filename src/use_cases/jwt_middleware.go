package use_cases

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header")
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// TODO: Change this
			return []byte("SECRETODA"), nil
		})

		// If parsing the token resulted in an error, return an unauthorized error
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Parse error: "+err.Error())
		}

		// If the token is valid, proceed with request
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("Token claims: ", claims)
			return next(c)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}
	}
}
