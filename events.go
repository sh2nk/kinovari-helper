package main

import (
	"context"
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v2/events"
)

// ParsedCommand - contains cmd action name with cmd arguments separatly
type ParsedCommand struct {
	Action string
	Args   []string
}

type Forward struct {
	OwnerID                 int   `json:"owner_id"`
	PeerID                  int   `json:"peer_id"`
	ConverstationMessageIDs []int `json:"conversation_message_ids"`
	IsReply                 bool  `json:"is_reply"`
}

// New message event handler
func onMessageNew(ctx context.Context, obj events.MessageNewObject) {
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
	switch cmd.Action {
	case "погода":
		var f Command
		f = Weather
		f = AddReply(f)
		f = SendMessage(f)

		if _, err := f(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}

	case "админ":
		var f Command
		f = Admin
		f = AddReply(f)
		f = SendMessage(f)

		if _, err := f(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'admin' command call: %v\n", err)
		}

	case "спасибо", "спс", "благодарю", "уважение", "респект", "хорош", "харош", "красава", "лучший", "лучшая", "хороша", "+":
		if obj.Message.ReplyMessage != nil {
			var f Command
			f = AddPhoto(BasicCommand, "img/thanks.jpg")
			f = AddReplyToSelected(f)
			f = SendMessage(f)

			if _, err := f(ctx, obj, cmd.Args); err != nil {
				log.Printf("Error occured during 'thanks' command call: %v\n", err)
			}
		}

	case "фу", "гавно", "говно", "дерьмо", "кал", "удали", "бан", "плох", "паршив", "паршива", "подводишь", "расстраиваешь", "опозорился", "опозорилась", "опростоволосился", "опростоволосилась":
		if obj.Message.ReplyMessage != nil {
			var c Command
			c = AddPhoto(BasicCommand, "img/oops.jpg")
			c = AddReplyToSelected(c)
			c = SendMessage(c)

			if _, err := c(ctx, obj, cmd.Args); err != nil {
				log.Printf("Error occured during 'oops' command call: %v\n", err)
			}
		}

	case "кнопки":
		kbd := NewKeyboard()
		kbd.AddRow()
		kbd.ButtonRows[0].AddButton("Кнопка 1", "positive")
		kbd.ButtonRows[0].AddButton("Кнопка 2", "primary")
		kbd.AddRow()
		kbd.ButtonRows[1].AddButton("Кнопка 4", "secondary")
		kbd.ButtonRows[1].AddButton("Кнопка 4", "negative")

		var f Command

		f = AddText(BasicCommand, "Тест кнопок")
		f = AddKeyboard(f, kbd)
		f = SendMessage(f)

		if _, err := f(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'buttons' command call: %v\n", err)
		}
	}
}
