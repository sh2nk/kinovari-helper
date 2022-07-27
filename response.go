package main

import (
	"context"
	"encoding/json"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// Command - basic command function type, used by response builder
type Command func(context.Context, events.MessageNewObject, []string) (*params.MessagesSendBuilder, error)

type Buttons struct {
}

// Response messages builder, wraps the command functions
func SendMessage(fn Command) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
		// Get result from inner function
		b, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		b.RandomID(int(randomInt32()))
		b.PeerID(obj.Message.PeerID)

		// Use forward field to reply if message from converstation
		if obj.Message.PeerID < 2000000000 {
			b.ReplyTo(obj.Message.ID)
		} else {
			f := Forward{
				PeerID:                  obj.Message.PeerID,
				ConverstationMessageIDs: []int{obj.Message.ConversationMessageID},
				IsReply:                 true,
			}
			bytes, err := json.Marshal(f)
			if err != nil {
				return nil, err
			}
			b.Forward(string(bytes))
		}

		// Sending response message
		_, err = VK.MessagesSend(b.Params)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
}

func AddPhoto(fn Command, path string) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
		// Get result from inner function
		b, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
}

func AddButtons(fn Command, btn string) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
		// Get result from inner function
		b, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
}
