package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/SevereCloud/vksdk/v2/events"
)

// Sends weather forecast
func Weather(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	if len(args) > 0 {
		req := Request{
			City:  args[0],
			AppID: WeatherToken,
			Units: "metric",
			Lang:  "ru",
		}

		resp, err := http.Get(req.MakeWeather())
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var data WeatherCurrentData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", err
		}

		t, err := template.ParseFiles("templates/weatherCurrent.txt")
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err = t.Execute(&buf, data); err != nil {
			return "", err
		}

		return buf.String(), nil
	} else {
		return "Не указан населенный пункт", nil
	}
}

// Sends user privilegies info
func Admin(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	if isConverstationAdmin(obj.Message.PeerID, obj.Message.FromID) {
		return "Пользователь админ в чате", nil
	} else {
		return "Пользователь простой работяга", nil
	}
}
