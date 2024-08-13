# most of the API refs are not mine, thanks to https://github.com/mukulhase/WebWhatsapp-Wrapper
import os
from .whatsmeow import (
    new_whatsapp_client_wrapper,
    connect_wrapper,
    disconnect_wrapper,
    message_thread_wrapper,
    send_message_wrapper,
    send_image_wrapper,
    send_video_wrapper,
    send_audio_wrapper,
    send_document_wrapper,
)
import ctypes
import json


class WhatsApp:
    """
    The main whatsapp handler
    """

    def __init__(
        self,
        phone_number: str = "",
        media_path: str = "",
        machine: str = "mac",
        browser: str = "safari",
        on_event=None,
        on_disconnect=None,
    ):
        """
        Import the compiled whatsmeow golang package, and setup basic client and database.
        Auto run based on any database (login and chat info database), hence a user phone number are declared.
        If there is no user login assigned yet, assign a new client.
        Put the database in current file whereever this class instances are imported. database/client.db
        :param phone_number: User phone number. in the Whatsmeow golang are called client.
        :param media_path: A directory to save all the media received
        :param machine: OS login info (showed on the whatsapp app)
        :param browser: Browser login info (showed on the whatsapp app)
        :param on_event: Function to call on event
        :param on_disconnect: Function to call on disconnect
        """

        self.user_name = None
        self.machine = machine
        self.browser = browser
        self.wapi_functions = browser
        self.connected = None

        if media_path:
            if not os.path.exists(media_path):
                os.makedirs(media_path)
            for subdir in ["images", "audios", "videos", "documents", "stickers"]:
                full_media_path = media_path + "/" + subdir
                if not os.path.exists(full_media_path):
                    os.makedirs(full_media_path)

        def on_event_json(s):
            try:
                s = s.decode()
            except:
                pass
            try:
                s = json.loads(s)
            except:
                pass

            on_event(s)

        CMPFUNC_NONE_STR = ctypes.CFUNCTYPE(None, ctypes.c_char_p)
        CMPFUNC_NONE = ctypes.CFUNCTYPE(None)

        self.C_ON_EVENT = (
            CMPFUNC_NONE_STR(on_event_json)
            if callable(on_event)
            else ctypes.cast(None, CMPFUNC_NONE_STR)
        )
        self.C_ON_DISCONNECT = (
            CMPFUNC_NONE(on_disconnect)
            if callable(on_disconnect)
            else ctypes.cast(None, CMPFUNC_NONE)
        )

        self.c_WhatsAppClientId = new_whatsapp_client_wrapper(
            phone_number.encode(),
            media_path.encode(),
            self.C_ON_DISCONNECT,
            self.C_ON_EVENT,
        )

    def connect(self):
        """
        Connects the whatsapp client to whatsapp servers. This method SHOULD be called before any other.
        """
        connect_wrapper(self.c_WhatsAppClientId)

    def disconnect(self):
        """
        Disconnects the whatsapp client to whatsapp servers.
        """
        disconnect_wrapper(self.c_WhatsAppClientId)

    def runMessageThread(self):
        """
        Checks for queued events and call on_event on new events.
        """
        message_thread_wrapper(self.c_WhatsAppClientId)

    def sendMessage(self, phone: str, message: str, group: bool = False):
        """
        Sends a text message
        :param phone: The phone number or group number to send the message.
        :param message: The message to send
        :param group: Send the message to a group ?
        :return: Function success or not
        """
        ret = send_message_wrapper(
            self.c_WhatsAppClientId, phone.encode(), message.encode(), group
        )
        return ret == 1

    def sendImage(
        self, phone: str, image_path: str, caption: str = "", group: bool = False
    ):
        """
        Sends a image message
        :param phone: The phone number or group number to send the message.
        :param image_path: The path to the image to send
        :param caption: The caption for the image
        :param group: Send the message to a group ?
        :return: Function success or not
        """
        ret = send_image_wrapper(
            self.c_WhatsAppClientId,
            phone.encode(),
            image_path.encode(),
            caption.encode(),
            group,
        )
        return ret == 1

    def sendVideo(
        self, phone: str, video_path: str, caption: str = "", group: bool = False
    ):
        """
        Sends a video message
        :param phone: The phone number or group number to send the message.
        :param video_path: The path to the video to send
        :param caption: The caption for the video
        :param group: Send the message to a group ?
        return: Function success or not
        """
        ret = send_video_wrapper(
            self.c_WhatsAppClientId,
            phone.encode(),
            video_path.encode(),
            caption.encode(),
            group,
        )
        return ret == 1

    def sendAudio(self, phone: str, audio_path: str, group: bool = False):
        raise NotImplementedError
        return send_audio_wrapper(
            self.c_WhatsAppClientId, phone.encode(), audio_path.encode(), group
        )

    def sendDocument(
        self, phone: str, document_path: str, caption: str, group: bool = False
    ):
        """
        Sends a document message
        :param phone: The phone number or group number to send the message.
        :param document_path: The path to the document to send
        :param caption: The caption for the document
        :param group: Send the message to a group ?
        return: Function success or not
        """
        return send_document_wrapper(
            self.c_WhatsAppClientId,
            phone.encode(),
            document_path.encode(),
            caption.encode(),
            group,
        )

    # -- unimplemented

    def get_all_chats(self):
        raise NotImplementedError
        return []

    def get_all_chat_ids(self):
        raise NotImplementedError
        return []

    def get_unread_messages_in_chat(
        self, id, include_me=False, include_notifications=False
    ):
        raise NotImplementedError

        # get unread messages
        # return them
        unread = []
        return unread

    def get_contacts(self):
        raise NotImplementedError

        return []

    def chat_send_seen(self, chat_id):
        raise NotImplementedError

        return self.wapi_functions.sendSeen(chat_id)

    def check_number_status(self, number_id) -> bool:
        raise NotImplementedError

        return True

    def subscribe_new_messages(self, observer):
        raise NotImplementedError

    def unsubscribe_new_messages(self, observer):
        raise NotImplementedError

    def is_connected(self) -> bool:
        raise NotImplementedError
        # return self.wapi_functions.isConnected()
        return True


if __name__ == "__main__":
    client = WhatsApp()
    message = "Hello World!"
    phone = "6283139000000"
    client.sendMessage(message=message, phone=phone)
