package main

import (
	"context"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/events"
)

// Sends weather forecast
func weather(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	return fmt.Sprintf("вызвана команда погода с аргументами: %#v", args), nil
}

// Sends user privilegies info
func admin(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
	if isConverstationAdmin(obj.Message.PeerID, obj.Message.FromID) {
		return "пользователь админ в чате.", nil
	} else {
		return "пользователь простой работяга.", nil
	}
}
