package controllers

import (
	"enchantech-codex/src/use_cases"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"net/http"
)

type FetchArticlesController struct {
	BaseController
	fetchArticlesUseCase *use_cases.FetchArticlesUseCase
}

func NewFetchArticlesController(fetchArticlesUseCase *use_cases.FetchArticlesUseCase) *FetchArticlesController {
	return &FetchArticlesController{
		BaseController: BaseController{
			Path:   "/api/articles",
			Method: http.MethodGet,
		},
		fetchArticlesUseCase: fetchArticlesUseCase,
	}
}

func (c *FetchArticlesController) Handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := new(use_cases.FetchArticlesUseCaseInput)

		decoder := schema.NewDecoder()
		err := decoder.Decode(input, ctx.QueryParams())
		if err != nil {
			return ctx.String(400, "Bad request")
		}

		fmt.Println(input)

		articles, err := c.fetchArticlesUseCase.Execute(*input)
		if err != nil {
			return ctx.String(400, "Bad request")
		}
		return ctx.JSON(200, articles)
		// TODO create DTOs to adjust the fields returned to the client
	}
}
