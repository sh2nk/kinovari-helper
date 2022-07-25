package main

import (
	"fmt"
	"net/url"
)

var (
	WeatherEndpoint  = "https://api.openweathermap.org/data/2.5/weather?"
	ForecastEndpoint = "https://api.openweathermap.org/data/2.5/forecast?"
)

type Coords struct {
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}

type WeatherInfo struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp        float32 `json:"temp"`
	FeelsLike   float32 `json:"feels_like"`
	TempMin     float32 `json:"temp_min"`
	TempMax     float32 `json:"temp_max"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
	SeaLevel    int     `json:"sea_level"`
	GroundLevel int     `json:"grnd_level"`
	TempKF      float32 `json:"temp_kf"`
}

type Wind struct {
	Speed     float32 `json:"speed"`
	Direction int     `json:"deg"`
	Gust      float32 `json:"gust"`
}

type Clouds struct {
	All int `json:"all"`
}

type Precipitation struct {
	OneHour    float32 `json:"1h"`
	ThreeHours float32 `json:"3h"`
}

type Sys struct {
	Type      int     `json:"type"`
	ID        int     `json:"id"`
	Message   float32 `json:"message"`
	Country   string  `json:"country"`
	Sunrise   int     `json:"sunrise"`
	Sunset    int     `json:"sunset"`
	PartOfDay string  `json:"pod"`
}

type WeatherCurrentData struct {
	Coords     Coords        `json:"coord"`
	Weather    []WeatherInfo `json:"weather"`
	Base       string        `json:"base"`
	Main       Main          `json:"main"`
	Visibility int           `json:"visibility"`
	Wind       Wind          `json:"wind"`
	Clouds     Clouds        `json:"clouds"`
	Rain       Precipitation `json:"rain"`
	Snow       Precipitation `json:"snow"`
	Date       int           `json:"dt"`
	Sys        Sys           `json:"sys"`
	Timezone   int           `json:"timezone"`
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Cod        int           `json:"cod"`
}

type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Coords     Coords `json:"coord"`
	Country    string `json:"country"`
	Population int    `json:"population"`
	Timezone   int    `json:"timezone"`
	Sunrise    int    `json:"sunrise"`
	Sunset     int    `json:"sunset"`
}

type WeatherItem struct {
	DT         int           `json:"dt"`
	Main       Main          `json:"main"`
	Weather    []WeatherInfo `json:"weather"`
	Clouds     Clouds        `json:"clouds"`
	Wind       Wind          `json:"wind"`
	Visibility int           `json:"visibility"`
	Pop        float32       `json:"pop"`
	Rain       Precipitation `json:"rain"`
	Snow       Precipitation `json:"snow"`
	Sys        Sys           `json:"sys"`
	Date       string        `json:"dt_txt"`
}

type WeatherForecast struct {
	Cod     int           `json:"cod"`
	Message int           `json:"message"`
	Count   int           `json:"cnt"`
	List    []WeatherItem `json:"list"`
	City    City          `json:"city"`
}

type Request struct {
	City  string
	AppID string
	Mode  string
	Units string
	Lang  string
	Count string
}

func (w Request) create() string {
	form := url.Values{}
	form.Add("q", w.City)
	form.Add("appid", w.AppID)
	form.Add("mode", w.Mode)
	form.Add("units", w.Units)
	form.Add("lang", w.Lang)

	return form.Encode()
}

func (w Request) MakeWeather() string {
	return WeatherEndpoint + w.create()
}

func (w Request) MakeForecast() string {
	return ForecastEndpoint + w.create()
}

func (w WeatherCurrentData) GetIcon() string {
	switch w.Weather[0].Icon {
	case "01d":
		return "☀️"
	case "01n":
		return "🌚"
	case "02d", "02n":
		return "🌤️"
	case "03d", "03n":
		return "☁️"
	case "04d", "04n":
		return "🌥️"
	case "09d", "09n":
		return "🌧️"
	case "10d", "10n":
		return "🌦️"
	case "11d", "11n":
		return "⛈️"
	case "13d", "13n":
		return "❄️"
	case "50d", "50n":
		return "🌫️"
	default:
		return "☀️"
	}
}

func (w WeatherCurrentData) GetTemp() string {
	if w.Main.Temp >= 0 {
		return fmt.Sprintf("+%.1f°", w.Main.Temp)
	} else {
		return fmt.Sprintf("-%.1f°", w.Main.Temp)
	}
}

func (w WeatherCurrentData) GetFeelsLike() string {
	if w.Main.FeelsLike >= 0 {
		return fmt.Sprintf("+%.1f°", w.Main.FeelsLike)
	} else {
		return fmt.Sprintf("-%.1f°", w.Main.FeelsLike)
	}
}

func (w WeatherCurrentData) GetPressure() float32 {
	return float32(w.Main.Pressure) * 0.750062
}
