package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func homeHandler(c echo.Context) error {

	scope := "read_station read_thermostat write_thermostat"
	state := "1234"
	baseURL := "https://api.netatmo.com/oauth2/authorize"

	clientIdParam := fmt.Sprintf("client_id=%s", clientID)
	redirectUriParam := fmt.Sprintf("redirect_uri=%s", url.QueryEscape(authRedirectURL))
	scopeParam := fmt.Sprintf("scope=%s", url.QueryEscape(scope))
	stateParam := fmt.Sprintf("state=%s", state)

	myURL := fmt.Sprintf("%s?%s&%s&%s&%s", baseURL, clientIdParam, redirectUriParam, scopeParam, stateParam)

	data := map[string]interface{}{
		"Title": "Welcome Page",
		"Url":   myURL,
	}
	return c.Render(http.StatusOK, "home.html", data)
}

func getCurrentWeather(c echo.Context) error {

	if err := ensureToken(); err != nil {
		logError(err, "ensuring access token")
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
	}

	weatherInfo, err := getWeather()
	if err != nil {
		logError(err, "getting current weather")
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}

	return c.JSONPretty(http.StatusOK, weatherInfo, " ")
}

func authRedirect(c echo.Context) error {

	returnedCode := c.QueryParam("code")
	if returnedCode == "" {
		// No code means this wasn't reached via Netatmo's redirect (or the
		// code was stripped). Re-start the flow from "/".
		if e := c.QueryParam("error"); e != "" {
			return c.String(http.StatusBadRequest, "Netatmo returned an authorization error: "+e+" "+c.QueryParam("error_description"))
		}
		return c.String(http.StatusBadRequest, "No authorization code present. Start the login again from /")
	}

	authUrl := "https://api.netatmo.com/oauth2/token"
	// Scope MUST match the scope requested in homeHandler's authorize URL,
	// or Netatmo rejects the code exchange.
	scope := "read_station read_thermostat write_thermostat"

	log.Printf("auth_redirect: exchanging code (len=%d) with redirect_uri=%q", len(returnedCode), authRedirectURL)

	postData := url.Values{}
	postData.Set("grant_type", "authorization_code")
	postData.Set("client_id", clientID)
	postData.Set("client_secret", clientSecret)
	postData.Set("code", returnedCode)
	postData.Set("redirect_uri", authRedirectURL)
	postData.Set("scope", scope)

	resp, err := http.PostForm(authUrl, postData)
	if err != nil {
		logError(err, "auth token POST request")
		return c.String(http.StatusBadGateway, "Error contacting Netatmo: "+err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logError(err, "reading auth response body")
		return c.String(http.StatusBadGateway, "Error reading Netatmo response: "+err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Netatmo returned HTTP %d: %s", resp.StatusCode, string(body))
		if strings.Contains(string(body), "invalid_grant") {
			msg += fmt.Sprintf("\n\ninvalid_grant means Netatmo rejected the code. Common causes:\n"+
				"  1. The code was already used or expired — start over at / and don't reload this page.\n"+
				"  2. redirect_uri mismatch — this exchange used %q, which must EXACTLY match the\n"+
				"     URI registered in your Netatmo dev console and the host you logged in from.",
				authRedirectURL)
		}
		logError(fmt.Errorf("%s", string(body)), "auth code exchange rejected")
		return c.String(http.StatusBadGateway, msg)
	}

	var authReturnData AuthReturn
	if err := json.Unmarshal(body, &authReturnData); err != nil {
		logError(err, "unmarshalling auth response")
		return c.String(http.StatusBadGateway, "Error parsing Netatmo response: "+err.Error())
	}
	if authReturnData.AccessToken == "" {
		return c.String(http.StatusBadGateway, "Netatmo response had no access token: "+string(body))
	}

	tokens.set(authReturnData.AccessToken, authReturnData.RefreshToken, authReturnData.ExpiresIn)

	return c.String(http.StatusOK, "Got tokens, should be good to go now")
}
