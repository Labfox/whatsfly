package main

// #include "wapp.h"
// #include <string.h>
// #include <stdlib.h>
// #include <stdint.h>
import "C"

import (
	"net/http"
	"os"

	// "os/signal"
	// "syscall"
	"context"
	"fmt"
	"path/filepath"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite"

	// sqlite3 "github.com/mattn/go-sqlite3"

	"encoding/json"
	"mime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/enriquebris/goconcurrentqueue"
	"go.mau.fi/whatsmeow/appstate"
	"google.golang.org/protobuf/encoding/protojson"
)

// var log waLog.Logger

func getJid(user string, is_group bool) types.JID {
	server := types.DefaultUserServer
	if is_group {
		server = types.GroupServer
	}

	return types.JID{
		User:   user,
		Server: server,
	}
}

type WhatsAppClient struct {
	phoneNumber          string
	mediaPath            string
	fnDisconnectCallback C.ptr_to_pyfunc
	fnEventCallback      C.ptr_to_pyfunc_str
	wpClient             *whatsmeow.Client
	eventQueue           *goconcurrentqueue.FIFO
	runMessageThread     bool
	isLoggedIn           bool
	startupTime          int64
	historySyncID        int32
}

var handles []*WhatsAppClient

func NewWhatsAppClient(phoneNumber string, mediaPath string, fn_disconnect_callback C.ptr_to_pyfunc, fn_event_callback C.ptr_to_pyfunc_str) *WhatsAppClient {
	return &WhatsAppClient{
		phoneNumber:          phoneNumber,
		mediaPath:            mediaPath,
		fnDisconnectCallback: fn_disconnect_callback,
		fnEventCallback:      fn_event_callback,
		wpClient:             nil,
		eventQueue:           goconcurrentqueue.NewFIFO(),
		runMessageThread:     false,
		isLoggedIn:           false,
		startupTime:          time.Now().Unix(),
		historySyncID:        0,
	}
}

func (w *WhatsAppClient) Connect() {
	// Set the path for the database file
	dbPath := "whatsapp/wapp.db"

	// Set Browser
	store.DeviceProps.PlatformType = waProto.DeviceProps_SAFARI.Enum()
	store.DeviceProps.Os = proto.String("macOS") //"Mac OS 10"

	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		panic(err)
	}

	// Connect to the database
	container, err := sqlstore.New("sqlite", "file:"+dbPath+"?_pragma=foreign_keys(1)", waLog.Noop)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	client := whatsmeow.NewClient(deviceStore, waLog.Noop)

	client.AddEventHandler(w.handler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}

	outerLoop:
		for {
			select {
			case <-time.After(60 * time.Second):
				w.Disconnect(client)
				w.eventQueue.Enqueue("{\"eventType\":\"disconnect\", \"reason\": \"timeout\"}")
				return
			case evt, ok := <-qrChan:
				if !ok {
					break outerLoop
				}
				if evt.Event == "code" {
					if len(w.phoneNumber) > 0 {
						linkingCode, err := client.PairPhone(w.phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
						if err != nil {
							panic(err)
						}
						w.eventQueue.Enqueue("{\"eventType\":\"linkCode\", \"code\": \"" + linkingCode + "\"}")
					}else {
                        w.eventQueue.Enqueue("{\"eventType\":\"qrCode\", \"code\": \"" + evt.Code + "\"}")
					}
				} else {
				}
			}
		}
	} else {
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	}

	w.wpClient = client
}

func (w *WhatsAppClient) addEventToQueue(msg string) {
	w.eventQueue.Enqueue(msg)
}

func (w *WhatsAppClient) handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(w.wpClient.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := w.wpClient.SendPresence(types.PresenceAvailable)
			if err != nil {
				//log.Warnf("Failed to send available presence: %v", err)
			} else {
				w.addEventToQueue("{\"eventType\":\"AppStateSyncComplete\",\"name\":\"" + string(evt.Name) + "\"}")
				//log.Infof("Marked self as available")
			}
		}
	case *events.Connected:
		if len(w.wpClient.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := w.wpClient.SendPresence(types.PresenceAvailable)
		if err != nil {
			//log.Warnf("Failed to send available presence: %v", err)
		} else {
			w.addEventToQueue("{\"eventType\":\"Connected\"}")
			//log.Infof("Marked self as available")
		}
	case *events.PushNameSetting:
		if len(w.wpClient.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := w.wpClient.SendPresence(types.PresenceAvailable)
		if err != nil {
			//log.Warnf("Failed to send available presence: %v", err)
		} else {
			w.addEventToQueue("{\"eventType\":\"PushNameSetting\",\"timestamp\":" + strconv.FormatInt(evt.Timestamp.Unix(), 10) + ",\"action\":\"" + (*evt.Action.Name) + "\",\"fromFullSync\":" + strconv.FormatBool(evt.FromFullSync) + "}")
			//log.Infof("Marked self as available")
		}
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:

		var info string


		info += "{\"id\":\"" + evt.Info.ID + "\""
		info += ",\"messageSource\":\"" + evt.Info.MessageSource.SourceString() + "\""
		if evt.Info.Type != "" {
			info += ",\"type\":\"" + evt.Info.Type + "\""
		}
		info += ",\"pushName\":\"" + evt.Info.PushName + "\""
		info += ",\"timestamp\":" + strconv.FormatInt(evt.Info.Timestamp.Unix(), 10)
		if evt.Info.Category != "" {
			info += ",\"category\":\"" + evt.Info.Category + "\""
		}
		info += ",\"multicast\":" + strconv.FormatBool(evt.Info.Multicast)
		if evt.Info.MediaType != "" {
			info += ",\"mediaType\": \"" + evt.Info.MediaType + "\""
		}
		info += ",\"flags\":["


		var flags []string
		if evt.IsEphemeral {
			flags = append(flags, "\"ephemeral\"")
		}
		if evt.IsViewOnce {
			flags = append(flags, "\"viewOnce\"")
		}
		if evt.IsViewOnceV2 {
			flags = append(flags, "\"viewOnceV2\"")
		}
		if evt.IsDocumentWithCaption {
			flags = append(flags, "\"documentWithCaption\"")
		}
		if evt.IsEdit {
			flags = append(flags, "\"edit\"")
		}
		info += strings.Join(flags, ",")
		info += "]"

        if evt.Message.ImageMessage != nil || evt.Message.AudioMessage != nil || evt.Message.VideoMessage != nil || evt.Message.DocumentMessage != nil || evt.Message.StickerMessage != nil {
            if len(w.mediaPath) > 0 {
                var mimetype string
                var media_path_subdir string
                var data []byte
                var err error
                switch {
                case evt.Message.ImageMessage != nil:
                    mimetype = evt.Message.ImageMessage.GetMimetype()
                    data, err = w.wpClient.Download(evt.Message.ImageMessage)
                    media_path_subdir = "images"
                case evt.Message.AudioMessage != nil:
                    mimetype = evt.Message.AudioMessage.GetMimetype()
                    data, err = w.wpClient.Download(evt.Message.AudioMessage)
                    media_path_subdir = "audios"
                case evt.Message.VideoMessage != nil:
                    mimetype = evt.Message.VideoMessage.GetMimetype()
                    data, err = w.wpClient.Download(evt.Message.VideoMessage)
                    media_path_subdir = "videos"
                case evt.Message.DocumentMessage != nil:
                    mimetype = evt.Message.DocumentMessage.GetMimetype()
                    data, err = w.wpClient.Download(evt.Message.DocumentMessage)
                    media_path_subdir = "documents"
                case evt.Message.StickerMessage != nil:
                    mimetype = evt.Message.StickerMessage.GetMimetype()
                    data, err = w.wpClient.Download(evt.Message.StickerMessage)
                    media_path_subdir = "stickers"
                }



                if err != nil {
                    fmt.Printf("Failed to download media: %v", err)
                } else {
                    exts, _ := mime.ExtensionsByType(mimetype)
                    path := fmt.Sprintf("%s/%s/%s%s", w.mediaPath, media_path_subdir, evt.Info.ID, exts[0])

                    err = os.WriteFile(path, data, 0600)
                    if err != nil {
                        fmt.Printf("Failed to save media: %v", err)
                    } else {
                        info += ",\"filepath\":\"" + path + "\""
                    }
                }



            }
        }

		info += "}"



		var m, _ = protojson.Marshal(evt.Message)
		var message_info string = string(m)
		json_str := "{\"eventType\":\"Message\",\"info\":" + info + ",\"message\":" + message_info + "}"

		w.addEventToQueue(json_str)



	case *events.Receipt:
		if evt.Type == types.ReceiptTypeRead || evt.Type == types.ReceiptTypeReadSelf {
			json_str := "{\"eventType\":\"MessageRead\",\"messageIDs\":["

			messageIDsLen := len(evt.MessageIDs)
			for key, value := range evt.MessageIDs {
				json_str += "\"" + value + "\""
				if key < messageIDsLen-1 {
					json_str += ","
				}
			}
			json_str += "],\"sourceString\":\"" + evt.SourceString() + "\",\"timestamp\":" + strconv.FormatInt(evt.Timestamp.Unix(), 10) + "}"

			w.addEventToQueue(json_str)
			//log.Infof("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
		} else if evt.Type == types.ReceiptTypeDelivered {
			json_str := "{\"eventType\":\"MessageDelivered\",\"messageID\":\"" + evt.MessageIDs[0] + "\",\"sourceString\":\"" + evt.SourceString() + "\",\"timestamp\":" + strconv.FormatInt(evt.Timestamp.Unix(), 10) + "}"
			w.addEventToQueue(json_str)
			//log.Infof("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
		}
	case *events.Presence:
		var json_str string
		var online bool = !evt.Unavailable

		json_str += "{\"eventType\":\"Presence\",\"from\":\"" + evt.From.String() + "\",\"online\":" + strconv.FormatBool(online)

		if evt.Unavailable {
			if !evt.LastSeen.IsZero() {
				json_str += ",\"lastSeen\":" + strconv.FormatInt(evt.LastSeen.Unix(), 10)
			}
		}
		json_str += "}"
		w.addEventToQueue(json_str)

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

		w.addEventToQueue("{\"eventType\":\"HistorySync\",\"filename\":\"" + fileName + "\"}")
	case *events.AppState:
		//log.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
		var json_str string = "{\"eventType\":\"AppState\",\"index\":["
		var event_index_size int = len(evt.Index)
		for key, value := range evt.Index {
			json_str += "\"" + value + "\""
			if key < event_index_size-1 {
				json_str += ","
			}
		}
		var protobuf_json, _ = protojson.Marshal(evt.SyncActionValue)
		var protobuf_json_str string = string(protobuf_json)
		// json_str := "{\"eventType\":\"Message\",\"info\":"+info+",\"message\":"+message_info+"}"

		json_str += "],\"syncActionValue\":" + protobuf_json_str + "}"
		// json_str += "],\"syncActionValue\":"+evt.SyncActionValue.String()+"}"

		w.addEventToQueue(json_str)

	case *events.KeepAliveTimeout:
		//log.Debugf("Keepalive timeout event: %+v", evt)
		var json_str string = "{\"eventType\":\"KeepAliveTimeout\",\"errorCount\":" + strconv.FormatInt(int64(evt.ErrorCount), 10) + ",\"lastSuccess\":" + strconv.FormatInt(evt.LastSuccess.Unix(), 10) + "}"
		w.addEventToQueue(json_str)
	case *events.KeepAliveRestored:
		//log.Debugf("Keepalive restored")
		w.addEventToQueue("{\"eventType\":\"KeepAliveRestored\"}")
	case *events.Blocklist:
		w.addEventToQueue("{\"eventType\":\"Blocklist\"}")
		//log.Infof("Blocklist event: %+v", evt)
	default:
		// fmt.Println("Missing event")
		// fmt.Printf("I don't know about type %T!\n", evt)

	}
}

func (w *WhatsAppClient) MessageThread() {
	w.runMessageThread = true
	for {
		if w.wpClient != nil {
			var is_logged_in_now = w.wpClient.IsLoggedIn()

			if w.isLoggedIn != is_logged_in_now {
				w.isLoggedIn = is_logged_in_now

				w.addEventToQueue("{\"eventType\":\"isLoggedIn\",\"loggedIn\":" + strconv.FormatBool(w.isLoggedIn) + "}")
				if !w.isLoggedIn {
					w.Disconnect(nil)
				}
			}
		}

		for w.eventQueue.GetLen() > 0 {
			value, _ := w.eventQueue.Dequeue()

			if w.fnEventCallback != nil {
				var str_value = value.(string)
				var cstr = C.CString(str_value)

				defer C.free(unsafe.Pointer(cstr))
				C.call_c_func_str(w.fnEventCallback, cstr)

			}
		}

		if !w.runMessageThread {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (w *WhatsAppClient) Disconnect(c2 *whatsmeow.Client) {
	client := w.wpClient

	if c2 != nil {
		client = c2
	}

	if client != nil {
		client.Disconnect()
	}

	if w.fnDisconnectCallback != nil {
		C.call_c_func(w.fnDisconnectCallback)
	}

	w.runMessageThread = false
}

func (w *WhatsAppClient) SendMessage(number string, message *waProto.Message, is_group bool) int {
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

func (w *WhatsAppClient) SendImage(number string, imagePath string, caption string, is_group bool) int {
	numberObj := getJid(number, is_group)

	// type imageStruct struct {
	//     Phone       string
	//     Image       string
	//     Caption     string
	//     Id          string
	//     ContextInfo waProto.ContextInfo
	// }
	// Check if the client is connected

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	// var filedata []byte
	filedata, err := os.ReadFile(imagePath)
	if err != nil {
		return 1
	}

	var uploaded whatsmeow.UploadResponse
	uploaded, err = w.wpClient.Upload(context.Background(), filedata, whatsmeow.MediaImage)
	if err != nil {
		return 1
	}
	// "data:image/png;base64,\""

	messageObj := &waProto.Message{ImageMessage: &waProto.ImageMessage{
		Caption:       proto.String(caption),
		Mimetype:      proto.String(http.DetectContentType(filedata)),

		URL: &uploaded.URL,
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256: uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}
	_, err = w.wpClient.SendMessage(context.Background(), numberObj, messageObj)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SendVideo(number string, videoPath string, caption string, is_group bool) int {
    numberObj := getJid(number, is_group)

	// type imageStruct struct {
	//     Phone       string
	//     Image       string
	//     Caption     string
	//     Id          string
	//     ContextInfo waProto.ContextInfo
	// }
	// Check if the client is connected

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	// var filedata []byte
	filedata, err := os.ReadFile(videoPath)
	if err != nil {
		return 1
	}


	var uploaded whatsmeow.UploadResponse
	uploaded, err = w.wpClient.Upload(context.Background(), filedata, whatsmeow.MediaVideo)
	if err != nil {
		return 1
	}

	// "data:image/png;base64,\""

	messageObj := &waProto.Message{VideoMessage: &waProto.VideoMessage{
		Caption:       proto.String(caption),
		Mimetype:      proto.String(http.DetectContentType(filedata)),

		URL: &uploaded.URL,
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256: uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}


	_, err = w.wpClient.SendMessage(context.Background(), numberObj, messageObj)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SendAudio(number string, audioPath string, is_group bool) int {
    numberObj := getJid(number, is_group)

	// type imageStruct struct {
	//     Phone       string
	//     Image       string
	//     Caption     string
	//     Id          string
	//     ContextInfo waProto.ContextInfo
	// }
	// Check if the client is connected

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	// var filedata []byte
	filedata, err := os.ReadFile(audioPath)
	if err != nil {
		return 1
	}


	var uploaded whatsmeow.UploadResponse
	uploaded, err = w.wpClient.Upload(context.Background(), filedata, whatsmeow.MediaAudio)
	if err != nil {
		return 1
	}



	// "data:image/png;base64,\""

	messageObj := &waProto.Message{AudioMessage: &waProto.AudioMessage{
		Mimetype: proto.String("audio/ogg; codec=opus"),
		URL: &uploaded.URL,
		PTT: proto.Bool(true),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256: uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}



	_, err = w.wpClient.SendMessage(context.Background(), numberObj, messageObj)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SendDocument(number string, documentPath string, caption string, is_group bool) int {
    numberObj := getJid(number, is_group)

	// type imageStruct struct {
	//     Phone       string
	//     Image       string
	//     Caption     string
	//     Id          string
	//     ContextInfo waProto.ContextInfo
	// }
	// Check if the client is connected

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	// var filedata []byte
	filedata, err := os.ReadFile(documentPath)
	if err != nil {
		return 1
	}


	var uploaded whatsmeow.UploadResponse
	uploaded, err = w.wpClient.Upload(context.Background(), filedata, whatsmeow.MediaDocument)
	if err != nil {
		return 1
	}

	// "data:image/png;base64,\""

	messageObj := &waProto.Message{DocumentMessage: &waProto.DocumentMessage{
		Caption:       proto.String(caption),
		Mimetype:      proto.String(http.DetectContentType(filedata)),

		URL: &uploaded.URL,
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256: uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}


	_, err = w.wpClient.SendMessage(context.Background(), numberObj, messageObj)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) GetGroupInviteLink(group string, reset bool, returnid string) int {
    numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	link, err := w.wpClient.GetGroupInviteLink(numberObj, reset)
	w.addEventToQueue("{\"eventType\":\"groupInviteLink\",\"group\": \""+group+"\", \"link\":\"" + link + "\"}")
	w.addEventToQueue("{\"eventType\":\"methodReturn\",\"return\": \""+link+"\", \"callid\":\"" + returnid + "\"}")
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) JoinGroupWithInviteLink(link string) int {
	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	_, err := w.wpClient.JoinGroupWithLink(link)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupAnnounce(group string, announce bool) int {
    numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupAnnounce(numberObj, announce)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupLocked(group string, locked bool) int {
    numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupLocked(numberObj, locked)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupName(group string, name string) int {
    numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupName(numberObj, name)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) SetGroupTopic(group string, topic string) int {
    numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	err := w.wpClient.SetGroupTopic(numberObj, "", "", topic)
	if err != nil {
		return 1
	}
	return 0
}

func (w *WhatsAppClient) GetGroupInfo(group string, return_id string) int {
    numberObj := getJid(group, true)

	if !w.wpClient.IsConnected() {
		err := w.wpClient.Connect()
		if err != nil {
			return 1
		}
	}

	groupinfo, err := w.wpClient.GetGroupInfo(numberObj)
	if err != nil {
		return 1
	}

    b, err := json.Marshal(&groupinfo)
    if err != nil {
        return 1
    }

    w.addEventToQueue("{\"eventType\":\"methodReturn\",\"return\": "+string(b)+", \"callid\":\"" + return_id + "\"}")

	return 0
}



//export NewWhatsAppClientWrapper
func NewWhatsAppClientWrapper(c_phone_number *C.char, c_media_path *C.char, fn_disconnect_callback C.ptr_to_pyfunc, fn_event_callback C.ptr_to_pyfunc_str) C.int {
	phone_number := C.GoString(c_phone_number)
	media_path := C.GoString(c_media_path)

	w := NewWhatsAppClient(phone_number, media_path, fn_disconnect_callback, fn_event_callback)
	handles = append(handles, w)

	return C.int(len(handles) - 1)
}

//export ConnectWrapper
func ConnectWrapper(id C.int) {
	w := handles[int(id)]
	w.Connect()
}

//export DisconnectWrapper
func DisconnectWrapper(id C.int) {
	w := handles[int(id)]
	w.Disconnect(nil)
}

//export MessageThreadWrapper
func MessageThreadWrapper(id C.int) {
	w := handles[int(id)]
	w.MessageThread()
}

//export SendMessageProtobufWrapper
func SendMessageProtobufWrapper(id C.int, c_phone_number *C.char, c_message *C.char, c_is_group C.bool) C.int {
	phone_number := C.GoString(c_phone_number)

	message := &waProto.Message{}

    length := C.strlen(c_message)

    goBytes := C.GoBytes(unsafe.Pointer(c_message), C.int(length))

    proto.Unmarshal(goBytes, message)
	is_group := bool(c_is_group)

	w := handles[int(id)]
	return C.int(w.SendMessage(phone_number, message, is_group))
}

//export SendMessageWrapper
func SendMessageWrapper(id C.int, c_phone_number *C.char, c_message *C.char, c_is_group C.bool) C.int {
    message_str := C.GoString(c_message)

    message := &waProto.Message{
        Conversation: &message_str,
    }

    marshalled, _ := proto.Marshal(message)

    message_encoded := C.CString(string(marshalled))

    return SendMessageProtobufWrapper(id, c_phone_number, message_encoded, c_is_group)
}

//export SendImageWrapper
func SendImageWrapper(id C.int, c_phone_number *C.char, c_image_path *C.char, c_caption *C.char, c_is_group C.bool) C.int {
	phone_number := C.GoString(c_phone_number)
	image_path := C.GoString(c_image_path)
	caption := C.GoString(c_caption)
	is_group := bool(c_is_group)

	w := handles[int(id)]
	return C.int(w.SendImage(phone_number, image_path, caption, is_group))
}

//export SendVideoWrapper
func SendVideoWrapper(id C.int, c_phone_number *C.char, c_video_path *C.char, c_caption *C.char, c_is_group C.bool) C.int {
	phone_number := C.GoString(c_phone_number)
	video_path := C.GoString(c_video_path)
	caption := C.GoString(c_caption)
	is_group := bool(c_is_group)

	w := handles[int(id)]
	return C.int(w.SendVideo(phone_number, video_path, caption, is_group))
}

//export SendAudioWrapper
func SendAudioWrapper(id C.int, c_phone_number *C.char, c_audio_path *C.char, c_is_group C.bool) C.int {
	phone_number := C.GoString(c_phone_number)
	audio_path := C.GoString(c_audio_path)
	is_group := bool(c_is_group)

	w := handles[int(id)]
	return C.int(w.SendAudio(phone_number, audio_path, is_group))
}

//export SendDocumentWrapper
func SendDocumentWrapper(id C.int, c_phone_number *C.char, c_document_path *C.char, c_caption *C.char, c_is_group C.bool) C.int {
	phone_number := C.GoString(c_phone_number)
	document_path := C.GoString(c_document_path)
	caption := C.GoString(c_caption)
	is_group := bool(c_is_group)

	w := handles[int(id)]
	return C.int(w.SendDocument(phone_number, document_path, caption, is_group))
}

//export GetGroupInviteLinkWrapper
func GetGroupInviteLinkWrapper(id C.int, c_jid *C.char, c_reset C.bool, c_return_id *C.char) C.int {
	jid := C.GoString(c_jid)
	reset := bool(c_reset)
	return_id := C.GoString(c_return_id)

	w := handles[int(id)]

	return C.int(w.GetGroupInviteLink(jid, reset, return_id))
}

//export JoinGroupWithInviteLinkWrapper
func JoinGroupWithInviteLinkWrapper(id C.int, c_link *C.char) C.int {
	link := C.GoString(c_link)

	w := handles[int(id)]

	return C.int(w.JoinGroupWithInviteLink(link))
}

//export SetGroupAnnounceWrapper
func SetGroupAnnounceWrapper(id C.int, c_jid *C.char, c_announce C.bool) C.int {
	jid := C.GoString(c_jid)
	announce := bool(c_announce)

	w := handles[int(id)]

	return C.int(w.SetGroupAnnounce(jid, announce))
}

//export SetGroupLockedWrapper
func SetGroupLockedWrapper(id C.int, c_jid *C.char, c_locked C.bool) C.int {
	jid := C.GoString(c_jid)
	locked := bool(c_locked)

	w := handles[int(id)]

	return C.int(w.SetGroupLocked(jid, locked))
}

//export SetGroupNameWrapper
func SetGroupNameWrapper(id C.int, c_jid *C.char, c_name *C.char) C.int {
	jid := C.GoString(c_jid)
	name := C.GoString(c_name)

	w := handles[int(id)]

	return C.int(w.SetGroupName(jid, name))
}

//export SetGroupTopicWrapper
func SetGroupTopicWrapper(id C.int, c_jid *C.char, c_topic *C.char) C.int {
	jid := C.GoString(c_jid)
	topic := C.GoString(c_topic)

	w := handles[int(id)]

	return C.int(w.SetGroupTopic(jid, topic))
}

//export GetGroupInfoWrapper
func GetGroupInfoWrapper(id C.int, c_jid *C.char, c_return_id *C.char) C.int {
	jid := C.GoString(c_jid)
	return_id := C.GoString(c_return_id)

	w := handles[int(id)]

	return C.int(w.GetGroupInfo(jid, return_id))
}


//export Version
func Version() C.int {
	return C.int(012)
}


func main() {
}
