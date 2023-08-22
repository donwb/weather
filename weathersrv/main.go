package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

var homeID string
var deviceID string
var btoken string
var rtoken string
var clientID string
var clientSecret string

func main() {

	setupEnvVars()

	e := echo.New()
	e.GET("/", homeHandler)
	e.GET("/current", getCurrentWeather)

	// Start!
	e.Logger.Fatal(e.Start(":1323"))
}

func setupEnvVars() {
	homeID = os.Getenv("HOMEID")
	deviceID = os.Getenv("DEVICEID")
	btoken = os.Getenv("BTOKEN")
	rtoken = os.Getenv("RTOKEN")
	clientID = os.Getenv("CLIENTID")
	clientSecret = os.Getenv("CLIENTSECRET")

	fmt.Println("\n\n-------- ENVVARS ------------")
	fmt.Println("HomeID: ", homeID)
	fmt.Println("DeviceID:", deviceID)
	fmt.Println("Bearer:", btoken)
	fmt.Println("Refresh", rtoken)
	fmt.Println("Client:", clientID)
	fmt.Println("Secret", clientSecret)
}
