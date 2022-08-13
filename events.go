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
		var c Command
		c = Weather
		c = AddReply(c)
		c = SendMessage(c)

		if _, err := c(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}

	case "админ":
		var c Command
		c = Admin
		c = AddReply(c)
		c = SendMessage(c)

		if _, err := c(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'admin' command call: %v\n", err)
		}

	case "спасибо", "спс", "благодарю", "уважение", "респект", "хорош", "харош", "красава", "лучший", "лучшая", "хороша", "+":
		if obj.Message.ReplyMessage != nil {
			var c Command
			c = AddPhoto(BasicCommand, "img/thanks.jpg")
			c = AddReplyToSelected(c)
			c = SendMessage(c)

			if _, err := c(ctx, obj, cmd.Args); err != nil {
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

	//TODO: make proper buttons constructor
	case "кнопки":
		kbd := &Keyboard{
			OneTime: false,
			Inline:  true,
			Buttons: [][]Button{
				{
					Button{
						Action: Action{
							Type:    "callback",
							Payload: "",
							Label:   "Кнопка 1",
						},
						Color: "positive",
					},
					Button{
						Action: Action{
							Type:    "callback",
							Payload: "",
							Label:   "Кнопка 2",
						},
						Color: "negative",
					},
				},
				{
					Button{
						Action: Action{
							Type:    "callback",
							Payload: "",
							Label:   "Кнопка 3",
						},
						Color: "primary",
					},
				},
			},
		}

		var c Command

		c = AddText(BasicCommand, "Тест кнопок")
		с = AddKeyboard(c, kbd)
		c = SendMessage(c)

		if _, err := c(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'buttons' command call: %v\n", err)
		}
	}
}
