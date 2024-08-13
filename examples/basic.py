from whatsfly import WhatsApp
import time

"""

basic usages

"""


def my_event_callback(event_data):
    """
    simple event callback to listen to incoming event/messages.
    whenever this function is called, it will retrieve the current incoming event or messages.
    """
    print("Received event data:", event_data)


def listening_message(minutes):
    """Stream messages for 'minutes' duration"""
    end_time = time.time() + minutes * 60
    while time.time() < end_time:
        time.sleep(1)


if __name__ == "__main__":
    phone = "6283139750000"  # make sure to attach country code + phone number
    message = "Hello World!"

    whatsapp = WhatsApp(on_event=my_event_callback)

    whatsapp.connect()

    message_sent = whatsapp.sendMessage(phone=phone, message=message)

    listening_message(5)

    whatsapp.disconnect()
