package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path"
	"sync/atomic"

	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/proto/waSyncAction"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/encoding/protojson"
)

type MessageEventCorpse struct {
	ID                  types.MessageID
	MessageSource       string
	Type                string
	PushName            string
	Timestamp           int64
	Category            string
	Multicast           bool
	Ephemeral           bool
	ViewOnce            bool
	ViewOnceV2          bool
	DocumentWithCaption bool
	Edit                bool
	MediaType           string
	Filepath            string
}

func (w *WhatsAppClient) handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(w.wpClient.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := w.wpClient.SendPresence(context.Background(), types.PresenceAvailable)
			if err != nil {
				//log.Warnf("Failed to send available presence: %v", err)
			} else {
				w.addEventToQueue(OutGoingEvent{
					Type:         "AppStateSyncComplete",
					JIDConcerned: w.wpClient.Store.GetJID().ADString(),
					Content:      evt.Name,
				})
				//log.Infof("Marked self as available")
			}
		}
	case *events.Connected:
		if len(w.wpClient.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := w.wpClient.SendPresence(context.Background(), types.PresenceAvailable)
		if err != nil {
			//log.Warnf("Failed to send available presence: %v", err)
		} else {
			w.addEventToQueue(OutGoingEvent{
				Type:         "Connected",
				JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			})
			//log.Infof("Marked self as available")
		}
	case *events.PushNameSetting:
		if len(w.wpClient.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := w.wpClient.SendPresence(context.Background(), types.PresenceAvailable)
		if err != nil {
			//log.Warnf("Failed to send available presence: %v", err)
		} else {
			//log.Infof("Marked self as available")
			w.addEventToQueue(OutGoingEvent{
				Type:         "PushNameSetting",
				JIDConcerned: w.wpClient.Store.GetJID().ADString(),
				Content: struct {
					Timestamp    int64  `json:"timestamp"`
					Action       string `json:"action"`
					FromFullSync bool   `json:"fromFullSync"`
				}{
					Timestamp:    evt.Timestamp.Unix(),
					Action:       *evt.Action.Name,
					FromFullSync: evt.FromFullSync,
				},
			})
		}
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:
		if evt.Message.GetPollUpdateMessage() != nil {
			decrytedPollvote, err := w.wpClient.DecryptPollVote(context.Background(), evt)
			if err == nil {
				data, _ := json.Marshal(decrytedPollvote)
				w.addEventToQueue(OutGoingEvent{
					Type:         "DecryptedPollvote",
					JIDConcerned: w.wpClient.Store.GetJID().ADString(),
					Content:      data,
				})
				return
			}
		}

		// TODO: correct serialisation
		info := MessageEventCorpse{
			ID:                  evt.Info.ID,
			MessageSource:       evt.Info.MessageSource.SourceString(),
			Type:                evt.Info.Type,
			PushName:            evt.Info.PushName,
			Timestamp:           evt.Info.Timestamp.Unix(),
			Category:            evt.Info.Category,
			Multicast:           evt.Info.Multicast,
			MediaType:           evt.Info.MediaType,
			Ephemeral:           evt.IsEphemeral,
			ViewOnce:            evt.IsViewOnce,
			ViewOnceV2:          evt.IsViewOnceV2,
			DocumentWithCaption: evt.IsDocumentWithCaption,
			Edit:                evt.IsEdit,
		}

		if evt.Message.ImageMessage != nil || evt.Message.AudioMessage != nil || evt.Message.VideoMessage != nil || evt.Message.DocumentMessage != nil || evt.Message.StickerMessage != nil {
			if len(w.mediaPath) > 0 {
				var mimetype string
				var media_path_subdir string
				var data []byte
				var err error
				switch {
				case evt.Message.ImageMessage != nil:
					mimetype = evt.Message.ImageMessage.GetMimetype()
					data, err = w.wpClient.Download(context.Background(), evt.Message.ImageMessage)
					media_path_subdir = "images"
				case evt.Message.AudioMessage != nil:
					mimetype = evt.Message.AudioMessage.GetMimetype()
					data, err = w.wpClient.Download(context.Background(), evt.Message.AudioMessage)
					media_path_subdir = "audios"
				case evt.Message.VideoMessage != nil:
					mimetype = evt.Message.VideoMessage.GetMimetype()
					data, err = w.wpClient.Download(context.Background(), evt.Message.VideoMessage)
					media_path_subdir = "videos"
				case evt.Message.DocumentMessage != nil:
					mimetype = evt.Message.DocumentMessage.GetMimetype()
					data, err = w.wpClient.Download(context.Background(), evt.Message.DocumentMessage)
					media_path_subdir = "documents"
				case evt.Message.StickerMessage != nil:
					mimetype = evt.Message.StickerMessage.GetMimetype()
					data, err = w.wpClient.Download(context.Background(), evt.Message.StickerMessage)
					media_path_subdir = "stickers"
				}

				if err != nil {
					fmt.Printf("Failed to download media: %v", err)
				} else {
					exts, _ := mime.ExtensionsByType(mimetype)
					fpath := path.Join(w.mediaPath, media_path_subdir, evt.Info.ID, exts[0])

					err = os.WriteFile(fpath, data, 0600)
					if err != nil {
						fmt.Printf("Failed to save media: %v", err)
					} else {
						info.Filepath = fpath
						w.addEventToQueue(OutGoingEvent{
							Type:         "MediaDownloaded",
							JIDConcerned: w.wpClient.Store.GetJID().ADString(),
							Content: struct {
								Path        string
								MessageInfo MessageEventCorpse
							}{
								Path:        fpath,
								MessageInfo: info,
							},
						})
					}
				}

			}
		}

		var m, _ = protojson.Marshal(evt.Message)
		var message_info string = string(m)

		w.addEventToQueue(OutGoingEvent{
			Type:         "Message",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content: struct {
				Info    MessageEventCorpse
				Message string
			}{
				Info:    info,
				Message: message_info,
			},
		})
		w.addEventToQueue(OutGoingEvent{
			Type:         "MessageJson",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content:      evt,
		})

	case *events.Receipt:
		var evt_type string
		switch evt.Type {
		case types.ReceiptTypeDelivered:
			evt_type = "delivered"
		case types.ReceiptTypeHistorySync:
			evt_type = "historySync"
		case types.ReceiptTypeInactive:
			evt_type = "inactive"
		case types.ReceiptTypePeerMsg:
			evt_type = "peerMsg"
		case types.ReceiptTypePlayed:
			evt_type = "played"
		case types.ReceiptTypePlayedSelf:
			evt_type = "playedSelf"
		case types.ReceiptTypeRead:
			evt_type = "read"
		case types.ReceiptTypeReadSelf:
			evt_type = "readSelf"
		case types.ReceiptTypeRetry:
			evt_type = "retry"
		case types.ReceiptTypeSender:
			evt_type = "sender"
		case types.ReceiptTypeServerError:
			evt_type = "serverError"
		}
		w.addEventToQueue(OutGoingEvent{
			Type:         "MessageDelivered",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content: struct {
				MessageIDs   []string `json:"messageIDs"`
				SourceString string   `json:"sourceString"`
				Timestamp    int64    `json:"timestamp"`
				Type         string   `json:"type"`
			}{
				MessageIDs:   evt.MessageIDs,
				SourceString: evt.SourceString(),
				Timestamp:    evt.Timestamp.Unix(),
				Type:         evt_type,
			},
		})
	case *events.Presence:
		w.addEventToQueue(OutGoingEvent{
			Type:         "Presence",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content: struct {
				From     string `json:"from"`
				Online   bool   `json:"online"`
				LastSeen int64  `json:"lastSeen"`
			}{
				From:     evt.From.String(),
				Online:   !evt.Unavailable,
				LastSeen: evt.LastSeen.Unix(),
			},
		})

	case *events.HistorySync:
		id := atomic.AddInt32(&w.historySyncID, 1)
		fileName := fmt.Sprintf("history-%d-%d.json", w.startupTime, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			//log.Errorf("Failed to open file to write history sync: %v", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		if err != nil {
			//log.Errorf("Failed to write history sync: %v", err)
			return
		}
		//log.Infof("Wrote history sync to %s", fileName)
		_ = file.Close()

		w.addEventToQueue(OutGoingEvent{
			Type:         "HistorySync",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content:      fileName,
		})
	case *events.AppState:
		//log.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)

		// var protobuf_json, _ = protojson.Marshal(evt.SyncActionValue)

		w.addEventToQueue(OutGoingEvent{
			Type:         "AppState",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content: struct {
				Index           []string                      `json:"index"`
				SyncActionValue *waSyncAction.SyncActionValue `json:"syncActionValue"`
			}{
				Index:           evt.Index,
				SyncActionValue: evt.SyncActionValue,
			},
		})

	case *events.KeepAliveTimeout:
		//log.Debugf("Keepalive timeout event: %+v", evt)
		w.addEventToQueue(OutGoingEvent{
			Type:         "KeepAliveTimeout",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
			Content: struct {
				ErrorCount  int   `json:"errorCount"`
				LastSuccess int64 `json:"lastSuccess"`
			}{ErrorCount: evt.ErrorCount, LastSuccess: evt.LastSuccess.Unix()},
		})
	case *events.KeepAliveRestored:
		//log.Debugf("Keepalive restored")

		w.addEventToQueue(OutGoingEvent{
			Type:         "KeepAliveRestored",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		})
	case *events.Blocklist:
		w.addEventToQueue(OutGoingEvent{
			Type:         "Blocklist",
			JIDConcerned: w.wpClient.Store.GetJID().ADString(),
		})
		//log.Infof("Blocklist event: %+v", evt)
	default:
		// fmt.Println("Missing event")
		// fmt.Printf("I don't know about type %T!\n", evt)

	}
}
