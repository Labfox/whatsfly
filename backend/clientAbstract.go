package main

// #include "wapp.h"
// #include <string.h>
// #include <stdlib.h>
// #include <stdint.h>
import "C"

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	"github.com/enriquebris/goconcurrentqueue"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsAppClient struct {
	phoneNumber          string
	mediaPath            string
	fnDisconnectCallback C.ptr_to_pyfunc
	fnEventCallback      C.ptr_to_pyfunc_str
	wpClient             *whatsmeow.Client
	eventQueue           *EQ
	runMessageThread     bool
	isLoggedIn           bool
	startupTime          int64
	historySyncID        int32
	uploadsData          map[string]whatsmeow.UploadResponse
	uploadsDataMutex     sync.Mutex
}

func NewWhatsAppClient(phoneNumber string, mediaPath string, fn_disconnect_callback C.ptr_to_pyfunc, fn_event_callback C.ptr_to_pyfunc_str) *WhatsAppClient {
	return &WhatsAppClient{
		phoneNumber:          phoneNumber,
		mediaPath:            mediaPath,
		fnDisconnectCallback: fn_disconnect_callback,
		fnEventCallback:      fn_event_callback,
		wpClient:             nil,
		eventQueue: &EQ{
			backend: goconcurrentqueue.NewFIFO(),
		},
		runMessageThread: false,
		isLoggedIn:       false,
		startupTime:      time.Now().Unix(),
		historySyncID:    0,
		uploadsData:      make(map[string]whatsmeow.UploadResponse),
	}
}

func (w *WhatsAppClient) addEventToQueue(msg OutGoingEvent) {
	w.eventQueue.Enqueue(msg)
}

func (w *WhatsAppClient) MessageThread() {
	w.runMessageThread = true
	for w.runMessageThread {
		if w.wpClient != nil {
			if !w.wpClient.IsConnected() {
				w.wpClient.Connect()
			}
			var is_logged_in_now = w.wpClient.IsLoggedIn()

			if w.isLoggedIn != is_logged_in_now {
				w.isLoggedIn = is_logged_in_now

				w.addEventToQueue(OutGoingEvent{
					Type:         "isLoggedIn",
					JIDConcerned: w.wpClient.Store.GetJID().ADString(),
					Content:      w.isLoggedIn,
				})
				if !w.isLoggedIn {
					w.Disconnect(nil)
				}
			}
		}

		for w.eventQueue.backend.GetLen() > 0 {
			value, _ := w.eventQueue.backend.Dequeue()

			if w.fnEventCallback != nil {
				var str_value = value.(string)
				var cstr = C.CString(str_value)

				defer C.free(unsafe.Pointer(cstr))
				C.call_c_func_str(w.fnEventCallback, cstr)

			}
		}

		time.Sleep(100 * time.Millisecond)

		if !w.runMessageThread {
			break
		}

	}
}

func (w *WhatsAppClient) Connect(dbPath string) {
	// Set the path for the database file
	//dbPath := "whatsapp/wapp.db"

	// Set Browser
	store.DeviceProps.PlatformType = waCompanionReg.DeviceProps_ANDROID_PHONE.Enum()
	store.DeviceProps.Os = proto.String("Android") //"Mac OS 10"

	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		panic(err)
	}

	// Connect to the database
	container, err := sqlstore.New(context.Background(), "sqlite", "file:"+dbPath+"?_pragma=foreign_keys(1)", waLog.Noop)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
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
				w.eventQueue.Enqueue(OutGoingEvent{
					Type:    "disconnect",
					Content: "timeout",
				})
				return
			case evt, ok := <-qrChan:
				if !ok {
					break outerLoop
				}
				if evt.Event == "code" {
					if len(w.phoneNumber) > 0 {
						linkingCode, err := client.PairPhone(context.Background(), w.phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
						if err != nil {
							panic(err)
						}
						w.eventQueue.Enqueue(OutGoingEvent{
							Type:    "linkingCode",
							Content: linkingCode,
						})
					} else {

						w.eventQueue.Enqueue(OutGoingEvent{
							Type:    "qrCode",
							Content: evt.Code,
						})
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
