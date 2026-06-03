package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const netatmoURL = "https://api.netatmo.com/api/"
const netatmoOauthURL = "https://api.netatmo.com/oauth2/token"

// ensureToken refreshes the access token only when it's missing or near expiry,
// avoiding a refresh (and a refresh-token rotation) on every single request.
func ensureToken() error {
	if !tokens.expired() {
		return nil
	}
	return refreshToken()
}

func refreshToken() error {
	cur := tokens.snapshot()
	if cur.RefreshToken == "" {
		return fmt.Errorf("no refresh token available; visit / to authorize")
	}

	postData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {cur.RefreshToken},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	req, err := http.NewRequest("POST", netatmoOauthURL, strings.NewReader(postData.Encode()))
	if err != nil {
		return fmt.Errorf("building refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading refresh response: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh returned HTTP %d: %s", res.StatusCode, string(body))
	}

	var ar AuthRefresh
	if err := json.Unmarshal(body, &ar); err != nil {
		return fmt.Errorf("parsing refresh response: %w", err)
	}
	if ar.AccessToken == "" {
		return fmt.Errorf("refresh response had no access token: %s", string(body))
	}

	tokens.set(ar.AccessToken, ar.RefreshToken, ar.ExpiresIn)
	return nil
}

func getWeather() (CurrentWeatherInfo, error) {
	queryParams := url.Values{
		"device_id": {deviceID},
		"home_id":   {homeID},
	}
	urlWithParams := netatmoURL + "homestatus?" + queryParams.Encode()

	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		return CurrentWeatherInfo{}, fmt.Errorf("building homestatus request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokens.snapshot().AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CurrentWeatherInfo{}, fmt.Errorf("homestatus request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CurrentWeatherInfo{}, fmt.Errorf("reading homestatus response: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return CurrentWeatherInfo{}, fmt.Errorf("homestatus returned HTTP %d: %s", res.StatusCode, string(body))
	}

	var homeStatus HomeStatus
	if err := json.Unmarshal(body, &homeStatus); err != nil {
		return CurrentWeatherInfo{}, fmt.Errorf("parsing homestatus response: %w", err)
	}
	if homeStatus.Status != "ok" {
		return CurrentWeatherInfo{}, fmt.Errorf("netatmo status %q: %s", homeStatus.Status, string(body))
	}

	return extractWeather(homeStatus), nil
}

// extractWeather maps modules by their reported type rather than array
// position, so an offline/reordered/missing module degrades gracefully
// (zeroed field) instead of panicking on a bad index.
func extractWeather(hs HomeStatus) CurrentWeatherInfo {
	var info CurrentWeatherInfo
	for _, m := range hs.Body.Home.Modules {
		switch m.Type {
		case "NAMain": // indoor base station
			info.InsideTemp = celsiusToFahrenheit(m.Temperature)
			info.Co2 = m.Co2
		case "NAModule1": // outdoor temperature/humidity
			info.OutsideTemp = celsiusToFahrenheit(m.Temperature)
			info.Humidity = m.Humidity
		case "NAModule3": // rain gauge, reported in mm -> convert to inches
			rain := m.SumRain24
			if rain > 0 {
				rain = rain / 25.4
			}
			info.Rainfall = rain
		}
	}
	return info
}
