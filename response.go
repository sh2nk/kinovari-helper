package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/SevereCloud/vksdk/v2/api/params"
)

// Command - basic command function type, used by response builder
type Command func(context.Context, MessageObject, []string) (*Message, error)

type Message struct {
	Builder *params.MessagesSendBuilder
}

// Response messages builder, wraps the command functions
func SendMessage(fn Command) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		m.Builder.RandomID(int(randomInt32()))
		m.Builder.PeerID(obj.Message.PeerID)

		// Sending response message
		if _, err = VK.MessagesSend(m.Builder.Params); err != nil {
			return nil, err
		}

		return m, nil
	}
}

// Send response by editing selected message instead of sending new message
// Should work for both new_message and event_message types
func EditMessage(fn Command, cmid int) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		m.Builder.RandomID(int(randomInt32()))

		if obj.Message.PeerID > 0 {
			m.Builder.PeerID(obj.Message.PeerID)
		} else {
			m.Builder.PeerID(obj.PeerID)
		}

		m.Builder.Params["conversation_message_id"] = cmid

		// Editing message
		if _, err = VK.MessagesEdit(m.Builder.Params); err != nil {
			return nil, err
		}

		return m, nil
	}
}

func AddReply(fn Command) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		// Use forward field to reply if message from converstation
		if obj.Message.PeerID < 2000000000 {
			m.Builder.ReplyTo(obj.Message.ID)
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
			m.Builder.Forward(string(bytes))
		}

		return m, nil
	}
}

func AddReplyToSelected(fn Command) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		// Use forward field to reply if message from converstation
		if obj.Message.PeerID < 2000000000 {
			m.Builder.ReplyTo(obj.Message.ReplyMessage.ID)
		} else {
			f := Forward{
				PeerID:                  obj.Message.PeerID,
				ConverstationMessageIDs: []int{obj.Message.ReplyMessage.ConversationMessageID},
				IsReply:                 true,
			}
			bytes, err := json.Marshal(f)
			if err != nil {
				return nil, err
			}
			m.Builder.Forward(string(bytes))
		}

		return m, nil
	}
}

func AddPhoto(fn Command, path string) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		// Open image file
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		photo, err := VK.UploadMessagesPhoto(obj.Message.PeerID, file)
		if err != nil {
			return nil, err
		}

		// Attach photo to message
		m.Builder.Attachment(photo)

		return m, nil
	}
}

func AddKeyboard(fn Command, kbd *Keyboard) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		k, err := json.Marshal(kbd)
		if err != nil {
			return nil, err
		}

		m.Builder.Keyboard(string(k))

		return m, nil
	}
}

func AddText(fn Command, t string) Command {
	return func(ctx context.Context, obj MessageObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		m.Builder.Message(t)

		return m, nil
	}
}
