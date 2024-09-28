"""
importing c shared whatsmeow library based on your machine.
"""

from sys import platform
from platform import machine
import ctypes
import os

# Load the shared library
if platform == "darwin":
    file_ext = "-darwin-arm64.dylib" if machine() == "arm64" else "-darwin-amd64.dylib"
elif platform in ("win32", "cygwin"):
    file_ext = (
        "-windows-64.dll" if 8 == ctypes.sizeof(ctypes.c_voidp) else "-windows-32.dll"
    )
else:
    machine = machine()
    if machine == "aarch64":
        file_ext = "-linux-arm64.so"
    elif machine.startswith("i686"):
        file_ext = "-linux-686.so"
    elif machine.startswith("i386"):
        file_ext = "-linux-386.so"
    else:
        file_ext = "-linux-amd64.so"

root_dir = os.path.abspath(os.path.dirname(__file__))

lib = ctypes.CDLL(f"{root_dir}/dependencies/whatsmeow/whatsmeow{file_ext}")


new_whatsapp_client_wrapper = lib.NewWhatsAppClientWrapper
new_whatsapp_client_wrapper.argstype = [
    ctypes.c_char_p,
    ctypes.c_char_p,
    ctypes.CFUNCTYPE(None),
    ctypes.CFUNCTYPE(None, ctypes.c_char_p),
]
new_whatsapp_client_wrapper.restype = ctypes.c_int

connect_wrapper = lib.ConnectWrapper
connect_wrapper.argstype = [ctypes.c_int]

disconnect_wrapper = lib.DisconnectWrapper
disconnect_wrapper.argstype = [ctypes.c_int]

logged_in_wrapper = lib.LoggedInWrapper
logged_in_wrapper.argstype = [ctypes.c_int]

connected_wrapper = lib.ConnectedWrapper
connected_wrapper.argstype = [ctypes.c_int]

message_thread_wrapper = lib.MessageThreadWrapper
message_thread_wrapper.argstype = [ctypes.c_int]

send_message_protobuf_wrapper = lib.SendMessageProtobufWrapper
send_message_protobuf_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p,
    ctypes.c_bool,
]

send_message_with_upload_wrapper = lib.SendMessageWithUploadWrapper
send_message_with_upload_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p,
    ctypes.c_bool,
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p
]


get_group_invite_link_wrapper = lib.GetGroupInviteLinkWrapper
get_group_invite_link_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_bool,
    ctypes.c_char_p,
]

join_group_with_invite_link_wrapper = lib.JoinGroupWithInviteLinkWrapper
join_group_with_invite_link_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p
]

set_group_announce_wrapper = lib.SetGroupAnnounceWrapper
set_group_announce_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_bool
]

set_group_locked_wrapper = lib.SetGroupLockedWrapper
set_group_locked_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_bool
]

set_group_name_wrapper = lib.SetGroupNameWrapper
set_group_name_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p
]

set_group_topic_wrapper = lib.SetGroupTopicWrapper
set_group_topic_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p
]

get_group_info_wrapper = lib.GetGroupInfoWrapper
get_group_info_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p
]

upload_file_wrapper = lib.UploadFileWrapper
upload_file_wrapper.argstype = [
    ctypes.c_int,
    ctypes.c_char_p,
    ctypes.c_char_p,
    ctypes.c_char_p
]