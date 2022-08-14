package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// Command - basic command function type, used by response builder
type Command func(context.Context, events.MessageNewObject, []string) (*Message, error)

type Message struct {
	Builder *params.MessagesSendBuilder
}

type Action struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label"`
}

type Button struct {
	Action Action `json:"action"`
	Color  string `json:"color"`
}

type Keyboard struct {
	OneTime    bool        `json:"one_time"`
	Inline     bool        `json:"inline"`
	ButtonRows []ButtonRow `json:"buttons"`
}

type ButtonRow []Button

func NewKeyboard() *Keyboard {
	return &Keyboard{OneTime: false, Inline: true}
}

func (k *Keyboard) AddRow() {
	k.ButtonRows = append(k.ButtonRows, ButtonRow{})
}

func (br *ButtonRow) AddButton(label, color string, payload ...string) {
	var p string
	if len(payload) > 0 {
		p = payload[0]
	} else {
		p = ""
	}
	b := Button{
		Action: Action{
			Type:    "callback",
			Payload: p,
			Label:   label,
		},
		Color: color,
	}
	*br = append(*br, b)
}

// Response messages builder, wraps the command functions
func SendMessage(fn Command) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
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

func AddReply(fn Command) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
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
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
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
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
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
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		k, err := json.Marshal(kbd)
		if err != nil {
			return nil, err
		}

		log.Println(string(k))
		m.Builder.Keyboard(string(k))

		return m, nil
	}
}

func AddText(fn Command, t string) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		m.Builder.Message(t)

		return m, nil
	}
}
