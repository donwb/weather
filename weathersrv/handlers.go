package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

	fmt.Println("URL: ", myURL)
	data := map[string]interface{}{
		"Title": "Welcome Page",
		"Url":   myURL,
	}
	//return c.String(200, "nothing to see here.......")
	return c.Render(http.StatusOK, "home.html", data)

}

func getCurrentWeather(c echo.Context) error {

	fmt.Println("\n\n\nGetting current weather....")
	tokenRefreshed := refreshToken()
	//tokenRefreshed := false

	fmt.Println("Token refreshed?: ", tokenRefreshed)
	weatherInfo := getWeather()

	return c.JSONPretty(200, weatherInfo, " ")

}

func authRedirect(c echo.Context) error {

	//searchQuery := c.QueryParam("q")
	returnedState := c.QueryParam("state")
	returnedCode := c.QueryParam("code")
	fmt.Println("Returned State: ", returnedState)
	fmt.Println("Returned Code: ", returnedCode)

	authUrl := "https://api.netatmo.com/oauth2/token"
	// need to send those params to the netatmo url and get back the refresh_token, auth_token and expires_in

	grantType := "authorization_code"
	redirectURI := authRedirectURL
	scope := "read_station read_thermostat write_thermostat"

	postData := url.Values{}
	postData.Set("grant_type", grantType)
	postData.Set("client_id", clientID)
	postData.Set("client_secret", clientSecret)
	postData.Set("code", returnedCode)
	postData.Set("redirect_uri", redirectURI)
	postData.Set("scope", scope)

	fmt.Println(url.Values(postData).Encode())

	// Make the POST request
	resp, err := http.PostForm(authUrl, postData)
	checkError(err, "Error making POST request")

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err, "Error reading response body")

	var authReturnData AuthReturn

	err = json.Unmarshal(body, &authReturnData)
	checkError(err, "Error unmarshalling JSON")

	// debug out the tokens for a minute until i'm sure this works
	fmt.Println("Access Token: ", authReturnData.AccessToken)
	fmt.Println("Refresh Token: ", authReturnData.RefreshToken)

	btoken = authReturnData.AccessToken
	rtoken = authReturnData.RefreshToken

	return c.String(200, "Got tokens, should be good to go now")

}
