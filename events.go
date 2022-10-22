package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/forPelevin/gomoji"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CommandString struct {
	Weather  string `yaml:"weather"`
	Forecast string `yaml:"forecast"`
	Thanks   string `yaml:"thanks"`
	Oops     string `yaml:"oops"`
}

// ParsedCommand - contains cmd action name with cmd arguments separatly
type ParsedCommand struct {
	Action string   `json:"action"`
	Args   []string `json:"args"`
}

type Forward struct {
	OwnerID                 int   `json:"owner_id"`
	PeerID                  int   `json:"peer_id"`
	ConverstationMessageIDs []int `json:"conversation_message_ids"`
	IsReply                 bool  `json:"is_reply"`
}

type MessageObject struct {
	events.MessageNewObject
	events.MessageEventObject
}

// New message event handler
func OnMessageNew(ctx context.Context, obj events.MessageNewObject) {
	// Get user input
	text := strings.Split(strings.ToLower(obj.Message.Text), " ")

	// Check if input starts with emoji
	if gomoji.ContainsEmoji(text[0]) {
		text = append(text[:0], text[1:]...)
	}

	// Attempting to parse command from user input
	cmd := ParsedCommand{
		Action: text[0],
		Args:   text[1:],
	}

	// Some debug logging action here
	if Debug {
		log.Printf("%d: %s", obj.Message.FromID, obj.Message.Text)
		log.Printf("Possible action: %s", cmd.Action)
	}

	// Do some actions if right prefix is found
	switch {
	// Weather command
	case strings.Contains(CommandStrings.Weather, cmd.Action):
		var f Command

		city := cases.Title(language.English).String(strings.Join(cmd.Args, " "))

		kbd := NewKeyboard()
		kbd.AddRows(3)
		kbd.ButtonRows[0].AddButton(
			fmt.Sprintf("🌂 Прогноз %s", city),
			"secondary",
			fmt.Sprintf("{\"action\":\"weather\",\"args\":[\"%s\"]}", city),
		)

		f = Weather
		f = AddReply(f)
		f = AddKeyboard(f, kbd)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}

	// Forecast command
	case strings.Contains(CommandStrings.Forecast, cmd.Action):
		var f Command

		city := cases.Title(language.English).String(strings.Join(cmd.Args, " "))

		kbd := NewKeyboard()
		kbd.AddRows(3)
		kbd.ButtonRows[0].AddButton(
			fmt.Sprintf("🌡️ Погода %s", city),
			"secondary",
			fmt.Sprintf("{\"action\":\"weather\",\"args\":[\"%s\"]}", city),
		)

		f = Forecast
		f = AddReply(f)
		f = AddKeyboard(f, kbd)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
			log.Printf("Error occured during 'forecast' command call: %v\n", err)
		}

	// Admin check (deprecated)
	case strings.Contains("///админ", cmd.Action):
		var f Command
		f = Admin
		f = AddReply(f)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
			log.Printf("Error occured during 'admin' command call: %v\n", err)
		}

	// Thanks command
	// TODO: add proper social rating system
	case strings.Contains(CommandStrings.Thanks, cmd.Action):
		if obj.Message.ReplyMessage != nil {
			var f Command
			f = AddPhoto(BasicCommand, "img/thanks.jpg")
			f = AddReplyToSelected(f)
			f = SendMessage(f)

			if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
				log.Printf("Error occured during 'thanks' command call: %v\n", err)
			}
		}

	// Oops command
	// TODO: add proper social rating system
	case strings.Contains(CommandStrings.Oops, cmd.Action):
		if obj.Message.ReplyMessage != nil {
			var c Command
			c = AddPhoto(BasicCommand, "img/oops.jpg")
			c = AddReplyToSelected(c)
			c = SendMessage(c)

			if _, err := c(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
				log.Printf("Error occured during 'oops' command call: %v\n", err)
			}
		}
	}
}

func OnMessageEvent(ctx context.Context, obj events.MessageEventObject) {
	var cmd ParsedCommand
	err := json.Unmarshal(obj.Payload, &cmd)
	if err != nil {
		log.Printf("Could not unmarshal event payload json: %v", err)
	}

	switch cmd.Action {
	// Forecast command
	case "forecast":
		var f Command

		kbd := NewKeyboard()
		kbd.AddRows(3)
		kbd.ButtonRows[0].AddButton(
			fmt.Sprintf("🌡️ Погода %s", cmd.Args[0]),
			"secondary",
			fmt.Sprintf("{\"action\":\"weather\",\"args\":[\"%s\"]}", cmd.Args[0]),
		)

		f = Forecast
		f = AddKeyboard(f, kbd)
		f = EditMessage(f, obj.ConversationMessageID)

		if _, err := f(ctx, MessageObject{events.MessageNewObject{}, obj}, cmd.Args); err != nil {
			log.Printf("Error occured during 'forecast' command call: %v\n", err)
		}

	// Weather command
	case "weather":
		var f Command

		kbd := NewKeyboard()
		kbd.AddRows(3)
		kbd.ButtonRows[0].AddButton(
			fmt.Sprintf("🌂 Прогноз %s", cmd.Args[0]),
			"secondary",
			fmt.Sprintf("{\"action\":\"forecast\",\"args\":[\"%s\"]}", cmd.Args[0]),
		)

		f = Weather
		f = AddKeyboard(f, kbd)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{events.MessageNewObject{}, obj}, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}
	}
}
