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
	switch cmd.Action {
	case "погода":
		if _, err := SendMessage(AddReply(Weather))(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}
	case "админ":
		if _, err := SendMessage(AddReply(Admin))(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'admin' command call: %v\n", err)
		}
	case "спасибо", "спс", "благодарю", "уважение", "респект", "хорош", "харош", "+":
		if obj.Message.ReplyMessage != nil {
			if _, err := SendMessage(AddReplyToSelected(AddPhoto(Empty, "img/thanks.jpg")))(ctx, obj, cmd.Args); err != nil {
				log.Printf("Error occured during 'thanks' command call: %v\n", err)
			}
		}
	case "фу", "гавно", "говно", "дерьмо", "кал", "удали", "бан", "плох", "паршив", "подводишь", "расстраиваешь":
		if obj.Message.ReplyMessage != nil {
			if _, err := SendMessage(AddReplyToSelected(AddPhoto(Empty, "img/oops.jpg")))(ctx, obj, cmd.Args); err != nil {
				log.Printf("Error occured during 'thanks' command call: %v\n", err)
			}
		}
	}
}
