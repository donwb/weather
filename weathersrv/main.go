package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

var homeID string
var deviceID string
var btoken string

func main() {

	homeID = os.Getenv("HOMEID")
	deviceID = os.Getenv("DEVICEID")
	btoken = os.Getenv("BTOKEN")

	fmt.Println("USING KEY: ", homeID)
	fmt.Println("-------- ENVVARS ------------")

	e := echo.New()
	e.GET("/", homeHandler)
	e.GET("/current", getCurrentWeather)

	// Start!
	e.Logger.Fatal(e.Start(":1323"))
}
