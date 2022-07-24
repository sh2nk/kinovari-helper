package main

import (
	"context"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/events"
)

type WeatherRequest struct {
	Location string
	AppID    string
	Mode     string
	Units    string
	Lang     string
	Count    string
}

type Coords struct {
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}

type Weather struct {
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
	Weather    Weather       `json:"weather"`
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
	Weather    Weather       `json:"weather"`
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

// Sends weather forecast
func weather(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {

	return fmt.Sprintf("вызвана команда погода с аргументами: %#v", args), nil
}

// Sends user privilegies info
func admin(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	if isConverstationAdmin(obj.Message.PeerID, obj.Message.FromID) {
		return "Пользователь админ в чате", nil
	} else {
		return "Пользователь простой работяга", nil
	}
}
