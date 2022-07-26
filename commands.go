package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"github.com/SevereCloud/vksdk/v2/events"
)

// Sends weather forecast
func Weather(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	// Check user input and throw a error if zero args is passed
	if len(args) > 0 {
		// Create weather API request
		req := Request{
			City:  strings.Join(args, " "),
			AppID: WeatherToken,
			Units: "metric",
			Lang:  "ru",
		}

		// Sending request
		resp, err := http.Get(req.MakeWeather())
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// Parsing responce data
		var data WeatherCurrentData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", err
		}

		// Response status codes handler
		switch data.Cod {
		case "404":
			return makeWarningMessage("Географический объект не найден."), nil
		case "429":
			return makeWarningMessage("Превышен лимит использования API!\n&#12288;Повторите попытку позже!"), nil
		default:
			// Parse message template
			t, err := template.ParseFiles("templates/weatherCurrent.txt")
			if err != nil {
				return "", err
			}

			// Render template to buffer
			var buf bytes.Buffer
			if err = t.Execute(&buf, data); err != nil {
				return "", err
			}
			return buf.String(), nil
		}
	} else {
		return makeErrorMessage("Не указан населенный пункт."), nil
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
