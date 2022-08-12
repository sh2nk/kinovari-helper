package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// Command - basic command function type, used by response builder
type Command func(context.Context, events.MessageNewObject, []string) (*Message, error)

type Message struct {
	Builder *params.MessagesSendBuilder
}

// Photo - uploaded image info
type Photo struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
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

		// // Get upload server URL
		// pb := params.NewPhotosGetMessagesUploadServerBuilder()
		// pb.PeerID(obj.Message.PeerID)
		// uresp, err := VK.PhotosGetMessagesUploadServer(pb.Params)
		// if err != nil {
		// 	return nil, err
		// }

		// // Create multipart
		// var buf bytes.Buffer
		// writer := multipart.NewWriter(&buf)
		// part, err := writer.CreateFormFile("photo", file.Name())
		// if err != nil {
		// 	return nil, err
		// }
		// part.Write(fdata)
		// writer.Close()

		// // Upload image
		// presp, err := http.Post(uresp.UploadURL, writer.FormDataContentType(), &buf)
		// if err != nil {
		// 	return nil, err
		// }
		// defer presp.Body.Close()

		// // Parse uploaded photo data
		// var pdata Photo
		// if err := json.NewDecoder(presp.Body).Decode(&pdata); err != nil {
		// 	return nil, err
		// }

		// // Make image save request and get attachment data
		// spb := params.NewPhotosSaveMessagesPhotoBuilder()
		// spb.Server(pdata.Server)
		// spb.Photo(pdata.Photo)
		// spb.Hash(pdata.Hash)
		// sresp, err := VK.PhotosSaveMessagesPhoto(spb.Params)
		// if err != nil {
		// 	return nil, err
		// }

		// Attach photo to message
		m.Builder.Attachment(photo)

		return m, nil
	}
}

func AddKeyboard(fn Command, kbd Keyboard) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*Message, error) {
		// Get result from inner function
		m, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		return m, nil
	}
}
