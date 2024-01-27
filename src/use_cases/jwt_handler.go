package use_cases

import (
	"enchantech-codex/src/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
}

func LoginHandler(c echo.Context) error {
	var req loginRequest

	utils.ErrorPanicPrinter(c.Bind(&req), false)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// encoded token
	tokenString, err := token.SignedString([]byte("SECRETODA"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})

}
