package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const netatmoURL = "https://api.netatmo.com/api/"
const netatmoOauthURL = "https://api.netatmo.com/oauth2/token"

func refreshToken() bool {

	postData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {rtoken},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	fmt.Println("Post data: ", postData)

	encodedPostData := postData.Encode()
	req, err := http.NewRequest("POST", netatmoOauthURL, strings.NewReader(encodedPostData))
	checkError(err, "ERROR: Constructing HTTP request")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	checkError(err, "ERROR: Making request")

	defer res.Body.Close()

	resBody := res.Body
	resBytes, err := ioutil.ReadAll(resBody)
	checkError(err, "ERROR: reading request")

	//fmt.Println("Response: ", string(resBytes))

	var authRefreshStatus AuthRefresh
	err = json.Unmarshal(resBytes, &authRefreshStatus)
	checkError(err, "ERROR: unmarshalling request")

	fmt.Println("Returned token: ", authRefreshStatus.AccessToken)
	if btoken != authRefreshStatus.AccessToken {
		btoken = authRefreshStatus.AccessToken
		rtoken = authRefreshStatus.RefreshToken
		fmt.Println("Refreshed access token to: ", btoken)
		return true
	} else {
		fmt.Println("No change in access token, continuing.... ")
		return false
	}

}

func getWeather() CurrentWeatherInfo {

	queryParams := url.Values{
		"device_id": {deviceID},
		"home_id":   {homeID},
	}
	homeStationURL := netatmoURL + "homestatus"
	urlWithParams := homeStationURL + "?" + queryParams.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlWithParams, nil)
	checkError(err, "error creating request")

	fmt.Println("Bearer: ", btoken)
	req.Header.Set("Authorization", "Bearer "+btoken)

	res, err := client.Do(req)
	checkError(err, "err making request")

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	checkError(err, "error reading request")

	var homeStatus HomeStatus
	err = json.Unmarshal(resBody, &homeStatus)
	checkError(err, "error unmarshalling request")

	returnString := string(resBody)
	hasError := strings.Contains(returnString, "error")
	if hasError {
		fmt.Println("Error in response: ", returnString)
		emptyReturnTemp := &CurrentWeatherInfo{}
		return *emptyReturnTemp
	}
	insideTemp := homeStatus.Body.Home.Modules[0].Temperature
	outsideTemp := homeStatus.Body.Home.Modules[1].Temperature
	co2 := homeStatus.Body.Home.Modules[0].Co2

	// reported in milimeters by netatmo
	rain := homeStatus.Body.Home.Modules[2].SumRain24
	if rain > 0 {
		rain = rain / 25.4
	}

	humidity := homeStatus.Body.Home.Modules[1].Humidity

	returnTemp := &CurrentWeatherInfo{
		OutsideTemp: celsiusToFahrenheit(outsideTemp),
		InsideTemp:  celsiusToFahrenheit(insideTemp),
		Rainfall:    rain,
		Humidity:    humidity,
		Co2:         co2,
	}

	return *returnTemp

}
