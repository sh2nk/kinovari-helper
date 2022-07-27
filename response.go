package main

import (
	"context"
	"encoding/json"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// Command - basic command function type, used by response builder
type Command func(context.Context, events.MessageNewObject, []string) (string, error)

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
		b.RandomID(int(randomInt32()))

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
				return res, err
			}
			b.Forward(string(bytes))
		}

		b.PeerID(obj.Message.PeerID)

		// Sending response message
		_, err = VK.MessagesSend(b.Params)
		if err != nil {
			return res, err
		}

		return res, nil
	}
}
