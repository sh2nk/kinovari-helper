package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SevereCloud/vksdk/v2/events"
)

// Sends weather forecast
func Weather(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	req := Request{
		City:  args[1],
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

	return fmt.Sprintf("Ответ сервера:\n%v", data), nil
}

// Sends user privilegies info
func Admin(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	if isConverstationAdmin(obj.Message.PeerID, obj.Message.FromID) {
		return "Пользователь админ в чате", nil
	} else {
		return "Пользователь простой работяга", nil
	}
}
