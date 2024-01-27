package di

import (
	"enchantech-codex/src/core/di"
	"enchantech-codex/src/use_cases"
	"fmt"
	"github.com/labstack/echo/v4"
)

func StartServer(container *di.Container) error {
	echoInstance := container.GetEchoInstance()
	feedController := container.GetFeedController()
	fetchArticlesController := container.GetFetchArticlesController()

	echoInstance.Add(fetchArticlesController.Method, fetchArticlesController.Path, fetchArticlesController.Handle())

	println("Configuring server routes...")
	echoInstance.GET("/api/force", feedController.UpdateArticles)
	echoInstance.POST("/api/login", use_cases.LoginHandler)
	echoInstance.GET("/api/protected", func(c echo.Context) error {
		return c.String(200, "{\"test\": \"test\"}")
	}, use_cases.JWTMiddleware)

	echoInstance.GET("*", func(c echo.Context) error {
		file := c.Request().URL.Path

		err := c.File("./dist/enchantech-codex-view/browser" + file)

		if err != nil {
			err = c.File("./dist/enchantech-codex-view/browser" + file + "/index.html")
			if err != nil {
				return c.File("./dist/enchantech-codex-view/browser/index.html")

			}
		}
		return fmt.Errorf("NOT FOUND")
	})

	println("Starting server...")
	return echoInstance.Start(":11001")
}
