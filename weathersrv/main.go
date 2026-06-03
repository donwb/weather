package main

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/labstack/echo/v4"
)

var homeID string
var deviceID string
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
	clientID = os.Getenv("CLIENTID")
	clientSecret = os.Getenv("CLIENTSECRET")
	authRedirectURL = os.Getenv("AUTHREDIRECT")

	tokenFile := os.Getenv("TOKENFILE")
	if tokenFile == "" {
		tokenFile = "tokens.json"
	}
	tokens = newTokenStore(tokenFile)
	tokens.bootstrap(os.Getenv("BTOKEN"), os.Getenv("RTOKEN"))

	log.Println("-------- CONFIG ------------")
	log.Println("HomeID:       ", homeID)
	log.Println("DeviceID:     ", deviceID)
	log.Println("ClientID:     ", clientID)
	log.Println("Auth Redirect:", authRedirectURL)
	log.Println("Token file:   ", tokenFile)
	log.Println("Have refresh token:", tokens.snapshot().RefreshToken != "")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
