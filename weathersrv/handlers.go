package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

var homeStationURL string

func homeHandler(c echo.Context) error {

	return c.String(200, "nothing to see here.......")

}

func getCurrentWeather(c echo.Context) error {

	homeStationURL := "https://api.netatmo.com/api/homestatus"

	queryParams := url.Values{
		"device_id": {deviceID},
		"home_id":   {homeID},
	}

	urlWithParams := homeStationURL + "?" + queryParams.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlWithParams, nil)
	checkError(err, "error creating request")

	req.Header.Set("Authorization", "Bearer "+btoken)

	res, err := client.Do(req)
	checkError(err, "err making request")

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	checkError(err, "error reading request")

	var homeStatus HomeStatus
	err = json.Unmarshal(resBody, &homeStatus)
	checkError(err, "error unmarshalling request")

	insideTemp := homeStatus.Body.Home.Modules[0].Temperature
	outsideTemp := homeStatus.Body.Home.Modules[1].Temperature
	rain := homeStatus.Body.Home.Modules[2].Rain
	humidity := homeStatus.Body.Home.Modules[1].Humidity

	returnTemp := &CurrentWeatherInfo{
		OutsideTemp: celsiusToFahrenheit(outsideTemp),
		InsideTemp:  celsiusToFahrenheit(insideTemp),
		Rainfall:    rain,
		Humidity:    humidity,
	}
	return c.JSONPretty(200, returnTemp, " ")

}
