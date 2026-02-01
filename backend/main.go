package main

// #include "wapp.h"
// #include <string.h>
// #include <stdlib.h>
// #include <stdint.h>
import "C"

import (

	// "os/signal"
	// "syscall"
	"context"
	"strconv"

	_ "embed"

	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite"

	// sqlite3 "github.com/mattn/go-sqlite3"

	"unsafe"
)

// var log waLog.Logger

var handles []*WhatsAppClient

//export NewWhatsAppClientWrapper
func NewWhatsAppClientWrapper(c_phone_number *C.char, c_media_path *C.char, fn_disconnect_callback C.ptr_to_pyfunc, fn_event_callback C.ptr_to_pyfunc_str) C.int {
	phone_number := C.GoString(c_phone_number)
	media_path := C.GoString(c_media_path)

	w := NewWhatsAppClient(phone_number, media_path, fn_disconnect_callback, fn_event_callback)
	defer w.eventQueue.Recoverer()
	handles = append(handles, w)

	return C.int(len(handles) - 1)
}

//export ConnectWrapper
func ConnectWrapper(id C.int, c_dbpath *C.char) {
	dbPath := C.GoString(c_dbpath)

	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	w.Connect(dbPath)
}

//export DisconnectWrapper
func DisconnectWrapper(id C.int) {
	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	w.Disconnect(nil)
}

//export ConnectedWrapper
func ConnectedWrapper(id C.int) C.int {
	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	if w.wpClient.IsConnected() {
		return C.int(1)
	} else {
		return C.int(0)
	}
}

//export LoggedInWrapper
func LoggedInWrapper(id C.int) C.int {
	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	if w.wpClient.IsLoggedIn() {
		return C.int(1)
	} else {
		return C.int(0)
	}
}

//export MessageThreadWrapper
func MessageThreadWrapper(id C.int) {
	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	w.MessageThread()
}

//export SendMessageProtobufWrapper
func SendMessageProtobufWrapper(id C.int, c_phone_number *C.char, c_message *C.char, c_is_group C.bool) C.int {
	phone_number := C.GoString(c_phone_number)

	message := &waE2E.Message{}

	length := C.strlen(c_message)

	goBytes := C.GoBytes(unsafe.Pointer(c_message), C.int(length))

	proto.Unmarshal(goBytes, message)
	is_group := bool(c_is_group)

	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	return C.int(w.SendMessage(phone_number, message, is_group))
}

//export SendMessageWithUploadWrapper
func SendMessageWithUploadWrapper(id C.int, c_phone_number *C.char, c_message *C.char, c_is_group C.bool, c_upload_id *C.char, c_mimetype *C.char, c_kind *C.char, c_ispb C.bool, c_thumbnail_path *C.char) C.int {
	phone_number := C.GoString(c_phone_number)
	w := handles[int(id)]
	defer w.eventQueue.Recoverer()
	mimetype := C.GoString(c_mimetype)

	kind := C.GoString(c_kind)

	caption := ""

	message := &waE2E.Message{}

	is_pb := bool(c_ispb)
	if is_pb {
		length := C.strlen(c_message)

		goBytes := C.GoBytes(unsafe.Pointer(c_message), C.int(length))

		proto.Unmarshal(goBytes, message)
	} else {
		caption = C.GoString(c_message)
	}

	is_group := bool(c_is_group)

	thumbnail_path := ""
	if c_thumbnail_path != nil {
		thumbnail_path = C.GoString(c_thumbnail_path)
	}

	upload_id := C.GoString(c_upload_id)

	w.uploadsDataMutex.Lock()
	defer w.uploadsDataMutex.Unlock()
	defer delete(w.uploadsData, upload_id)

	injected := w.InjectMessageWithUploadData(message, w.uploadsData[upload_id], mimetype, kind, caption, thumbnail_path)

	return C.int(w.SendMessage(phone_number, injected, is_group))
}

//export GetGroupInviteLinkWrapper
func GetGroupInviteLinkWrapper(id C.int, c_jid *C.char, c_reset C.bool, c_return_id *C.char) C.int {
	jid := C.GoString(c_jid)
	reset := bool(c_reset)
	return_id := C.GoString(c_return_id)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.GetGroupInviteLink(jid, reset, return_id))
}

//export JoinGroupWithInviteLinkWrapper
func JoinGroupWithInviteLinkWrapper(id C.int, c_link *C.char) C.int {
	link := C.GoString(c_link)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.JoinGroupWithInviteLink(link))
}

//export SetGroupAnnounceWrapper
func SetGroupAnnounceWrapper(id C.int, c_jid *C.char, c_announce C.bool) C.int {
	jid := C.GoString(c_jid)
	announce := bool(c_announce)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.SetGroupAnnounce(jid, announce))
}

//export SetGroupLockedWrapper
func SetGroupLockedWrapper(id C.int, c_jid *C.char, c_locked C.bool) C.int {
	jid := C.GoString(c_jid)
	locked := bool(c_locked)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.SetGroupLocked(jid, locked))
}

//export SetGroupNameWrapper
func SetGroupNameWrapper(id C.int, c_jid *C.char, c_name *C.char) C.int {
	jid := C.GoString(c_jid)
	name := C.GoString(c_name)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.SetGroupName(jid, name))
}

//export SetGroupTopicWrapper
func SetGroupTopicWrapper(id C.int, c_jid *C.char, c_topic *C.char) C.int {
	jid := C.GoString(c_jid)
	topic := C.GoString(c_topic)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.SetGroupTopic(jid, topic))
}

//export GetGroupInfoWrapper
func GetGroupInfoWrapper(id C.int, c_jid *C.char, c_return_id *C.char) C.int {
	jid := C.GoString(c_jid)
	return_id := C.GoString(c_return_id)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.GetGroupInfo(jid, return_id))
}

//export UploadFileWrapper
func UploadFileWrapper(id C.int, c_path *C.char, c_kind *C.char, c_return_id *C.char) C.int {
	path := C.GoString(c_path)
	return_id := C.GoString(c_return_id)
	kind := C.GoString(c_kind)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	return C.int(w.UploadFile(path, kind, return_id))
}

//export SendReactionWrapper
func SendReactionWrapper(id C.int, c_jid *C.char, c_message_id *C.char, c_sender_jid *C.char, c_reaction *C.char, c_group C.bool) C.int {
	message_id := C.GoString(c_message_id)
	reaction := C.GoString(c_reaction)

	numberObj := getJid(C.GoString(c_jid), bool(c_group))
	senderJID := getJid(C.GoString(c_sender_jid), false)

	w := handles[int(id)]

	defer w.eventQueue.Recoverer()

	_, err := w.wpClient.SendMessage(context.Background(), numberObj, w.wpClient.BuildReaction(numberObj, senderJID, message_id, reaction))
	if err != nil {
		return 1
	}
	return 0
}

//go:embed version.txt
var version string

//export Version
func Version() C.int {

	k, err := strconv.Atoi(version)
	if err != nil {
		return C.int(5)
	}

	return C.int(k)
}

func main() {
}
