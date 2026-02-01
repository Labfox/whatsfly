package main

import (
	"context"
	"os"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (w *WhatsAppClient) SendMessage(number string, message *waE2E.Message, is_group bool) int {
	var numberObj types.JID = getJid(number, is_group)

	messageObj := message

	// Check if the client is connected
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	// for {
	//     if w.wpClient.IsLoggedIn() {
	//         fmt.Println("Logged in!")
	//         break
	//     }
	// }

	_, err := w.wpClient.SendMessage(context.Background(), numberObj, messageObj)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) UploadFile(path string, kind string, return_id string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	// var filedata []byte
	filedata, err := os.ReadFile(path)
	if err != nil {
		return 1
	}

	var mediakind whatsmeow.MediaType

	if kind == "image" {
		mediakind = whatsmeow.MediaImage
	}
	if kind == "video" {
		mediakind = whatsmeow.MediaVideo
	}
	if kind == "audio" {
		mediakind = whatsmeow.MediaAudio
	}
	if kind == "document" {
		mediakind = whatsmeow.MediaDocument
	}

	var uploaded whatsmeow.UploadResponse
	uploaded, err = w.wpClient.Upload(context.Background(), filedata, mediakind)
	if err != nil {
		return 1
	}

	// w.uploadsData = append(w.uploadsData, uploaded)
	// data_return_uuid := len(w.uploadsData) - 1
	// Get the id and set it to data_return_uuid

	w.uploadsDataMutex.Lock()
	w.uploadsData[return_id] = uploaded
	w.uploadsDataMutex.Unlock()

	//w.addEventToQueue("{\"eventType\":\"methodReturn\",\"return\": \"" + strconv.Itoa(data_return_uuid) + "\", \"callid\":\"" + return_id + "\"}")
	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID: return_id,
		},
	})

	return 0
}

func (w *WhatsAppClient) InjectMessageWithUploadData(originMessage *waE2E.Message, upload whatsmeow.UploadResponse, mimetype string, kind string, caption string, thumbnail_path string) *waE2E.Message {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			panic(err)
		}
	}

	thumbnail_data, _ := os.ReadFile(thumbnail_path)

	// var filedata []byte

	if kind == "image" {
		originMessage.ImageMessage = &waE2E.ImageMessage{}

		if caption != "" {
			originMessage.ImageMessage.Caption = proto.String(caption)
		}

		originMessage.ImageMessage.Mimetype = proto.String(mimetype)
		originMessage.ImageMessage.URL = &upload.URL
		originMessage.ImageMessage.DirectPath = proto.String(upload.DirectPath)
		originMessage.ImageMessage.MediaKey = upload.MediaKey
		originMessage.ImageMessage.FileEncSHA256 = upload.FileEncSHA256
		originMessage.ImageMessage.FileSHA256 = upload.FileSHA256
		originMessage.ImageMessage.FileLength = proto.Uint64(uint64(upload.FileLength))

		if thumbnail_data != nil {
			originMessage.ImageMessage.JPEGThumbnail = thumbnail_data
		}
	}
	if kind == "video" {
		originMessage.VideoMessage = &waE2E.VideoMessage{}
		if caption != "" {
			originMessage.VideoMessage.Caption = proto.String(caption)
		}
		originMessage.VideoMessage.Mimetype = proto.String(mimetype)
		originMessage.VideoMessage.URL = &upload.URL
		originMessage.VideoMessage.DirectPath = proto.String(upload.DirectPath)
		originMessage.VideoMessage.MediaKey = upload.MediaKey
		originMessage.VideoMessage.FileEncSHA256 = upload.FileEncSHA256
		originMessage.VideoMessage.FileSHA256 = upload.FileSHA256
		originMessage.VideoMessage.FileLength = proto.Uint64(uint64(upload.FileLength))
	}
	if kind == "audio" {
		originMessage.AudioMessage = &waE2E.AudioMessage{}
		originMessage.AudioMessage.Mimetype = proto.String(mimetype)
		originMessage.AudioMessage.URL = &upload.URL
		originMessage.AudioMessage.DirectPath = proto.String(upload.DirectPath)
		originMessage.AudioMessage.MediaKey = upload.MediaKey
		originMessage.AudioMessage.FileEncSHA256 = upload.FileEncSHA256
		originMessage.AudioMessage.FileSHA256 = upload.FileSHA256
		originMessage.AudioMessage.FileLength = proto.Uint64(uint64(upload.FileLength))
	}
	if kind == "document" {
		originMessage.DocumentMessage = &waE2E.DocumentMessage{}
		if caption != "" {
			originMessage.DocumentMessage.Caption = proto.String(caption)
		}
		originMessage.DocumentMessage.Mimetype = proto.String(mimetype)
		originMessage.DocumentMessage.URL = &upload.URL
		originMessage.DocumentMessage.DirectPath = proto.String(upload.DirectPath)
		originMessage.DocumentMessage.MediaKey = upload.MediaKey
		originMessage.DocumentMessage.FileEncSHA256 = upload.FileEncSHA256
		originMessage.DocumentMessage.FileSHA256 = upload.FileSHA256
		originMessage.DocumentMessage.FileLength = proto.Uint64(uint64(upload.FileLength))
	}

	return originMessage
}
