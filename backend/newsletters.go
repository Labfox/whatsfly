package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (w *WhatsAppClient) CreateNewsletter(name string, description string, picturePath string, return_id string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	params := whatsmeow.CreateNewsletterParams{
		Name:        name,
		Description: description,
	}

	var picture []byte
	if picturePath != "" {
		fmt.Println("Setting the newsletter picture is unsupported. See Labfox/whatsfly#363")
		var err error
		picture, err = os.ReadFile(picturePath)
		if err != nil {
			return 1
		}
		params.Picture = picture
	}

	metadata, err := w.wpClient.CreateNewsletter(context.Background(), params)
	if err != nil {
		return 1
	}

	b, err := json.Marshal(&metadata)
	if err != nil {
		return 1
	}

	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID:      return_id,
			Content: b,
		},
	})
	return 0
}

func (w *WhatsAppClient) GetNewsletterInfo(jid string, return_id string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	jidObj := getNewsletterJid(jid)
	metadata, err := w.wpClient.GetNewsletterInfo(context.Background(), jidObj)
	if err != nil {
		return 1
	}

	b, err := json.Marshal(metadata)
	if err != nil {
		return 1
	}

	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID:      return_id,
			Content: b,
		},
	})
	return 0
}

func (w *WhatsAppClient) GetNewsletterMessages(jid string, count int, before int, return_id string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	jidObj := getNewsletterJid(jid)
	params := &whatsmeow.GetNewsletterMessagesParams{
		Count:  count,
		Before: types.MessageServerID(before),
	}

	messages, err := w.wpClient.GetNewsletterMessages(context.Background(), jidObj, params)
	if err != nil {
		return 1
	}

	b, err := json.Marshal(messages)
	if err != nil {
		return 1
	}

	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID:      return_id,
			Content: b,
		},
	})
	return 0
}

func (w *WhatsAppClient) GetSubscribedNewsletters(return_id string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	newsletters, err := w.wpClient.GetSubscribedNewsletters(context.Background())
	if err != nil {
		return 1
	}

	b, err := json.Marshal(newsletters)
	if err != nil {
		return 1
	}

	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID:      return_id,
			Content: b,
		},
	})
	return 0
}

func (w *WhatsAppClient) UploadNewsletter(path string, kind string, return_id string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	filedata, err := os.ReadFile(path)
	if err != nil {
		return 1
	}

	var mediakind whatsmeow.MediaType
	if kind == "image" {
		mediakind = whatsmeow.MediaImage
	} else if kind == "video" {
		mediakind = whatsmeow.MediaVideo
	} else if kind == "audio" {
		mediakind = whatsmeow.MediaAudio
	} else if kind == "document" {
		mediakind = whatsmeow.MediaDocument
	}

	uploaded, err := w.wpClient.UploadNewsletter(context.Background(), filedata, mediakind)
	if err != nil {
		return 1
	}

	w.uploadsDataMutex.Lock()
	w.uploadsData[return_id] = uploaded
	w.uploadsDataMutex.Unlock()

	w.addEventToQueue(OutGoingEvent{
		Type:         "methodReturn",
		JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		Content: ReturnFunc{
			ID: return_id,
		},
	})
	return 0
}

func (w *WhatsAppClient) SendNewsletter(jid string, message *waE2E.Message, mediaHandle string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	jidObj := getNewsletterJid(jid)

	if mediaHandle != "" {
		_, err := w.wpClient.SendMessage(context.Background(), jidObj, message, whatsmeow.SendRequestExtra{MediaHandle: mediaHandle})
		if err != nil {
			return 1
		}
	} else {
		_, err := w.wpClient.SendMessage(context.Background(), jidObj, message)
		if err != nil {
			return 1
		}
	}
	return 0
}
