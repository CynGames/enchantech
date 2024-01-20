package di

import (
	"enchantech-codex/src/core/di"
	"fmt"
	"github.com/labstack/echo/v4"
)

func StartServer(container *di.Container) error {
	echoInstance := container.GetEchoInstance()
	feedController := container.GetFeedController()

	println("Configuring server routes...")
	//echoInstance.GET("/", feedController.GetArticles)
	echoInstance.GET("/force", feedController.UpdateArticles)
	echoInstance.POST("/update", feedController.UpdateArticles)

	echoInstance.GET("/api/test-json", func(c echo.Context) error {

		return c.String(200, "{\"test\": \"test\"}")
	})

	echoInstance.GET("*", func(c echo.Context) error {
		fairu := c.Request().URL.Path

		err := c.File("./dist/enchantech-codex-view/browser" + fairu)

		if err != nil {
			err = c.File("./dist/enchantech-codex-view/browser" + fairu + "/index.html")
			if err != nil {
				return c.File("./dist/enchantech-codex-view/browser/index.html")

			}
		}
		return fmt.Errorf("NOT FOUNDO")
	})

	println("Starting server...")
	return echoInstance.Start(":11001")
}
