package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"text/template"

	"github.com/SevereCloud/vksdk/v2/api/params"
)

// Sends weather forecast
func Weather(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
	m := new(Message)
	m.Builder = params.NewMessagesSendBuilder()

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
		case "401":
			return nil, errors.New("invalid weather API key")
		case "404":
			m.Builder.Message(makeWarningMessage("Географический объект не найден."))
			return m, nil
		case "429":
			m.Builder.Message(makeWarningMessage("Превышен лимит использования API!\n&#12288;Повторите попытку позже!"))
			return m, nil
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
			m.Builder.Message(buf.String())
			return m, nil
		}
	} else {
		m.Builder.Message(makeErrorMessage("Не указан населенный пункт."))
		return m, nil
	}
}

// Sends user privilegies info
func Admin(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
	m := new(Message)
	m.Builder = params.NewMessagesSendBuilder()

	if isConverstationAdmin(obj.Message.PeerID, obj.Message.FromID) {
		m.Builder.Message("Пользователь админ в чате")
		return m, nil
	} else {
		m.Builder.Message("Пользователь простой работяга")
		return m, nil
	}
}

func BasicCommand(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
	return &Message{Builder: params.NewMessagesSendBuilder()}, nil
}
