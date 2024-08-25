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
    get_group_invite_link_wrapper,
    join_group_with_invite_link_wrapper,
    set_group_announce_wrapper,
    set_group_locked_wrapper,
    set_group_name_wrapper,
    set_group_topic_wrapper
)
import ctypes
import json
import threading
import warnings
import functools
import qrcode

def deprecated(func):
    """This is a decorator which can be used to mark functions
    as deprecated. It will result in a warning being emitted
    when the function is used."""
    @functools.wraps(func)
    def new_func(*args, **kwargs):
        warnings.simplefilter('always', DeprecationWarning)  # turn off filter
        warnings.warn("Call to deprecated function {}.".format(func.__name__),
                      category=DeprecationWarning,
                      stacklevel=2)
        warnings.simplefilter('default', DeprecationWarning)  # reset filter
        return func(*args, **kwargs)
    return new_func


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
        print_qr_code=True
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
        self._messageThreadRunner = threading.Thread(target=self._messageThread)
        self._userEventHandlers = [on_event]
        self.print_qr_code = print_qr_code

        if media_path:
            if not os.path.exists(media_path):
                os.makedirs(media_path)
            for subdir in ["images", "audios", "videos", "documents", "stickers"]:
                full_media_path = media_path + "/" + subdir
                if not os.path.exists(full_media_path):
                    os.makedirs(full_media_path)


        CMPFUNC_NONE_STR = ctypes.CFUNCTYPE(None, ctypes.c_char_p)
        CMPFUNC_NONE = ctypes.CFUNCTYPE(None)

        self.C_ON_EVENT = (
            CMPFUNC_NONE_STR(self._handleMessage)
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

        self._messageThreadRunner.start()

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

    @deprecated
    def runMessageThread(self):
        """
        Checks for queued events and call on_event on new events.
        """
        print("This method does nothing anymore, it has been automatised")

    def _messageThread(self):
        """
        New method for runMessageThread
        """
        while True:
            message_thread_wrapper(self.c_WhatsAppClientId)

    def _handleMessage(self, message):
        try:
            message = message.decode()
        except:
            pass
        try:
            message = json.loads(message)
        except:
            pass

        match message["eventType"]:
            case "linkCode":
                if self.print_qr_code:
                    print(message["code"])
            case "qrCode":
                if self.print_qr_code:
                    print(message["code"])
                    qr = qrcode.QRCode()
                    qr.add_data(message["code"])
                    qr.print_ascii()


        for handler in self._userEventHandlers:
            handler(message)

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
        :return: Function success or not
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
        :return: Function success or not
        """
        return send_document_wrapper(
            self.c_WhatsAppClientId,
            phone.encode(),
            document_path.encode(),
            caption.encode(),
            group,
        )

    def getGroupInviteLink(
            self, group: str, reset: bool = False
    ):
        """
        Get invite link for group, sends it to message queue
        :param group: Group id
        :param reset: If true, resets the old link before generating the new one
        :return: Successfull
        """
        return get_group_invite_link_wrapper(
            self.c_WhatsAppClientId,
            group.encode(),
            reset,
        )

    def joinGroupWithInviteLink(self, code: str):
        """
        Joins a group with an invite link
        :param code: The link
        """
        return join_group_with_invite_link_wrapper(
            self.c_WhatsAppClientId,
            code.encode(),
        )

    def setGroupAnnounce(self, group, announce):
        """
        Set a group's announce mode (only admins can send message
        :param group: Group jid
        :param announce: Announce mode or not
        """
        return set_group_announce_wrapper(
            self.c_WhatsAppClientId,
            group.encode(),
            announce
        )

    def setGroupLocked(self, group, locked):
        """
            Set a group's lock mode (only admins can change settings)
            :param group: Group jid
            :param locked: Lock mode or not
        """
        return set_group_locked_wrapper(
            self.c_WhatsAppClientId,
            group.encode(),
            locked
        )

    def setGroupName(self, group, name):
        """
            Set a group's name
            :param group: Group jid
            :param name: Name
        """
        return set_group_name_wrapper(
            self.c_WhatsAppClientId,
            group.encode(),
            name.encode()
        )

    def setGroupTopic(self, group, topic):
        """
        Set a group's topic
        :param group: Group jid
        :param topic: Topic
        """
        return set_group_topic_wrapper(
            self.c_WhatsAppClientId,
            group.encode(),
            topic.encode()
        )

    # -- unimplemented




if __name__ == "__main__":
    client = WhatsApp()
    message = "Hello World!"
    phone = "6283139000000"
    client.sendMessage(message=message, phone=phone)
