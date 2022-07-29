package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

// Command - basic command function type, used by response builder
type Command func(context.Context, events.MessageNewObject, []string) (*params.MessagesSendBuilder, error)

// Photo - uploaded image info
type Photo struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

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
		if _, err = VK.MessagesSend(b.Params); err != nil {
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

		// Open image file
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		fdata, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		file.Close()

		// Get upload server URL
		pb := params.NewPhotosGetMessagesUploadServerBuilder()
		pb.PeerID(obj.Message.PeerID)
		uresp, err := VK.PhotosGetMessagesUploadServer(pb.Params)
		if err != nil {
			return nil, err
		}

		// Create multipart
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, err := writer.CreateFormFile("photo", file.Name())
		if err != nil {
			return nil, err
		}
		part.Write(fdata)
		writer.Close()

		// Upload image
		presp, err := http.Post(uresp.UploadURL, writer.FormDataContentType(), &buf)
		if err != nil {
			return nil, err
		}
		defer presp.Body.Close()

		// Parse uploaded photo data
		var pdata Photo
		if err := json.NewDecoder(presp.Body).Decode(&pdata); err != nil {
			return nil, err
		}

		// Make image save request and get attachment data
		spb := params.NewPhotosSaveMessagesPhotoBuilder()
		spb.Server(pdata.Server)
		spb.Photo(pdata.Photo)
		spb.Hash(pdata.Hash)
		sresp, err := VK.PhotosSaveMessagesPhoto(spb.Params)
		if err != nil {
			return nil, err
		}

		// Attach photo to message
		b.Attachment(sresp)

		return b, nil
	}
}

func AddButtons(fn Command, btn Buttons) Command {
	return func(ctx context.Context, obj events.MessageNewObject, args []string) (*params.MessagesSendBuilder, error) {
		// Get result from inner function
		b, err := fn(ctx, obj, args)
		if err != nil {
			return nil, err
		}

		return b, nil
	}
}
