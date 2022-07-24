package main

import (
	"context"
	"log"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// ParsedCommand - contains cmd action name with cmd arguments separatly
type ParsedCommand struct {
	Action string
	Args   []string
}

// Command - basic command function type, used by response builder
type Command func(context.Context, events.MessageNewObject, []string) (string, error)

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
		if _, err := MakeResponse(Weather)(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'weather' command call: %v\n", err)
		}
	case "админ":
		if _, err := MakeResponse(Admin)(ctx, obj, cmd.Args); err != nil {
			log.Printf("Error occured during 'admin' command call: %v\n", err)
		}
	}
}

// Response messages builder, wraps the command functions
func MakeResponse(fn Command) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (string, error) {
		// Get result from inner function
		res, err := fn(ctx, obj, args)
		if err != nil {
			return res, err
		}

		// Building response
		b := params.NewMessagesSendBuilder()
		b.Message(res)
		b.RandomID(0)
		b.PeerID(obj.Message.PeerID)

		// Sending response message
		_, err = VK.MessagesSend(b.Params)
		if err != nil {
			return res, err
		}

		return res, nil
	}
}
