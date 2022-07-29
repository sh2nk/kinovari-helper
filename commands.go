package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// Sends weather forecast
func Weather(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
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
			return nil, err
		}
		defer resp.Body.Close()

		// Parsing responce data
		var data WeatherCurrentData
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, err
		}

		// Response status codes handler
		switch data.Cod {
		case "404":
			b.Message(makeWarningMessage("Географический объект не найден."))
			return b, nil
		case "429":
			b.Message(makeWarningMessage("Превышен лимит использования API!\n&#12288;Повторите попытку позже!"))
			return b, nil
		default:
			// Parse message template
			t, err := template.ParseFiles("templates/weatherCurrent.txt")
			if err != nil {
				return nil, err
			}

			// Render template to buffer
			var buf bytes.Buffer
			if err = t.Execute(&buf, data); err != nil {
				return nil, err
			}
			b.Message(buf.String())
			return b, nil
		}
	} else {
		b.Message(makeErrorMessage("Не указан населенный пункт."))
		return b, nil
	}
}

// Sends user privilegies info
func Admin(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
	b := params.NewMessagesSendBuilder()
	if isConverstationAdmin(obj.Message.PeerID, obj.Message.FromID) {
		b.Message("Пользователь админ в чате")
		return b, nil
	} else {
		b.Message("Пользователь простой работяга")
		return b, nil
	}
}

func Thanks(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
	return params.NewMessagesSendBuilder(), nil
}

func TestPhoto(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
	return params.NewMessagesSendBuilder(), nil
}
