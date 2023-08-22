package main

type CurrentWeatherInfo struct {
	InsideTemp  int `json:"insideTemp"`
	OutsideTemp int `json:"outsideTemp"`
	Rainfall    int `json:"rainfall"`
	Humidity    int `json:"humidity"`
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
				Rain             int     `json:"rain,omitempty"`
				SumRain1         int     `json:"sum_rain_1,omitempty"`
				SumRain24        int     `json:"sum_rain_24,omitempty"`
			} `json:"modules"`
		} `json:"home"`
	} `json:"body"`
}
