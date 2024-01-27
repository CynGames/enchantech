package controllers

import "github.com/labstack/echo/v4"

type BaseController struct {
	Method string
	Path   string
}

type Controller interface {
	Handle() echo.HandlerFunc
}
