package main

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

var homeID string
var deviceID string
var btoken string
var rtoken string
var clientID string
var clientSecret string
var authRedirectURL string

func main() {

	setupEnvVars()

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t
	e.GET("/", homeHandler)
	e.GET("/current", getCurrentWeather)
	e.GET("/auth_redirect", authRedirect)

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
	authRedirectURL = os.Getenv("AUTHREDIRECT")

	// Clearning these out for now, will remove from envvars later
	btoken = ""
	rtoken = ""

	fmt.Println("\n\n-------- ENVVARS ------------")
	fmt.Println("HomeID: ", homeID)
	fmt.Println("DeviceID:", deviceID)
	fmt.Println("Bearer:", btoken)
	fmt.Println("Refresh", rtoken)
	fmt.Println("Client:", clientID)
	fmt.Println("Secret", clientSecret)
	fmt.Println("Auth Redirect: ", authRedirectURL)
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
