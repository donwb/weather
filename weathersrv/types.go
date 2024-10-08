package main

type CurrentWeatherInfo struct {
	InsideTemp  int     `json:"insideTemp"`
	OutsideTemp int     `json:"outsideTemp"`
	Rainfall    float64 `json:"rainfall"`
	Humidity    int     `json:"humidity"`
	Co2         int     `json:"co2"`
}

type HomeStatus struct {
	Status     string `json:"status"`
	TimeServer int    `json:"time_server"`
	Body       struct {
		Home struct {
			ID      string `json:"id"`
			Modules []struct {
				ID               string  `json:"id"`
				Type             string  `json:"type"`
				FirmwareRevision int     `json:"firmware_revision"`
				WifiState        string  `json:"wifi_state,omitempty"`
				WifiStrength     int     `json:"wifi_strength,omitempty"`
				Ts               int     `json:"ts"`
				Temperature      float64 `json:"temperature,omitempty"`
				Co2              int     `json:"co2,omitempty"`
				Humidity         int     `json:"humidity,omitempty"`
				Noise            int     `json:"noise,omitempty"`
				Pressure         float64 `json:"pressure,omitempty"`
				AbsolutePressure float64 `json:"absolute_pressure,omitempty"`
				BatteryState     string  `json:"battery_state,omitempty"`
				BatteryLevel     int     `json:"battery_level,omitempty"`
				RfState          string  `json:"rf_state,omitempty"`
				RfStrength       int     `json:"rf_strength,omitempty"`
				LastSeen         int     `json:"last_seen,omitempty"`
				Reachable        bool    `json:"reachable,omitempty"`
				Bridge           string  `json:"bridge,omitempty"`
				Rain             float64 `json:"rain,omitempty"`
				SumRain1         float64 `json:"sum_rain_1,omitempty"`
				SumRain24        float64 `json:"sum_rain_24,omitempty"`
			} `json:"modules"`
		} `json:"home"`
	} `json:"body"`
}

type AuthRefresh struct {
	Scope        []string `json:"scope"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	ExpireIn     int      `json:"expire_in"`
}

type AuthReturn struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	ExpireIn     int      `json:"expire_in"`
	Scope        []string `json:"scope"`
}
