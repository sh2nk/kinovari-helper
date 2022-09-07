package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v2/events"
)

type CommandString struct {
	Weather string `yaml:"weather"`
	Thanks  string `yaml:"thanks"`
	Oops    string `yaml:"oops"`
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
	// Some debug logging action here
	if Debug {
		log.Printf("%d: %s", obj.Message.FromID, obj.Message.Text)
	}

	// Attempting to parse command from user input
	text := strings.Split(strings.ToLower(obj.Message.Text), " ")
	cmd := ParsedCommand{
		Action: text[0],
		Args:   text[1:],
	}

	// Do some actions if right prefix is found
	// TODO: export command keystrings to separate place
	switch {
	case strings.Contains(CommandStrings.Weather, cmd.Action):
		var f Command
		f = Weather
		f = AddReply(f)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}

	case strings.Contains("админ", cmd.Action):
		var f Command
		f = Admin
		f = AddReply(f)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
			log.Printf("Error occured during 'admin' command call: %v\n", err)
		}

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

	case strings.Contains("кнопки", cmd.Action):
		kbd := NewKeyboard()
		kbd.AddRow()
		kbd.ButtonRows[0].AddButton("Кнопка 1", "positive", "{\"action\": \"forecast\"}")
		kbd.ButtonRows[0].AddButton("Кнопка 2", "primary")
		kbd.AddRow()
		kbd.ButtonRows[1].AddButton("Кнопка 4", "secondary")
		kbd.ButtonRows[1].AddButton("Кнопка 4", "negative")

		var f Command
		f = AddText(BasicCommand, "Тест кнопок")
		f = AddKeyboard(f, kbd)
		f = SendMessage(f)

		if _, err := f(ctx, MessageObject{obj, events.MessageEventObject{}}, cmd.Args); err != nil {
			log.Printf("Error occured during 'buttons' command call: %v\n", err)
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
	case "forecast":
		var f Command
		f = AddText(BasicCommand, "Тест кнопок")
		f = EditMessage(f, obj.ConversationMessageID)

		if _, err := f(ctx, MessageObject{events.MessageNewObject{}, obj}, cmd.Args); err != nil {
			log.Printf("Error occured during 'forecast' command call: %v\n", err)
		}
	}
}
