package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func homeHandler(c echo.Context) error {

	return c.String(200, "nothing to see here.......")

}

func getCurrentWeather(c echo.Context) error {

	tokenRefreshed := refreshToken()
	fmt.Println("Token refreshed?: ", tokenRefreshed)
	weatherInfo := getWeather()

	return c.JSONPretty(200, weatherInfo, " ")

}
