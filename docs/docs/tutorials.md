# Tutorials

This page contains step-by-step guides for common use cases with WhatsFly.

## 1. Simple Echo Bot

An Echo Bot responds to every text message it receives by sending the same text back to the sender.

### Code

```python
from whatsfly import WhatsApp
import json
import time

def echo_callback(whatsapp, event):
    # Check if the event is an incoming message
    if event.get("eventType") == "Message":
        # The 'Message' field contains a JSON string of the message protobuf
        message_content = json.loads(event.get("Message", "{}"))
        info = event.get("Info", {})

        # Check if it's a conversation (text message)
        if "conversation" in message_content:
            text = message_content["conversation"]
            # Extract the sender's JID (phone number)
            # MessageSource usually looks like '1234567890@s.whatsapp.net/AD...'
            sender_jid = info.get("MessageSource", "").split("@")[0]

            print(f"Received message: '{text}' from {sender_jid}")

            # Echo the message back
            whatsapp.sendMessage(sender_jid, f"Echo: {text}")

if __name__ == "__main__":
    # Initialize WhatsApp client
    whatsapp = WhatsApp()

    # Register callback (can also be passed in constructor)
    whatsapp._userEventHandlers = [echo_callback]

    # Connect and wait for QR code scan
    whatsapp.connect()

    print("Echo bot is running. Press Ctrl+C to stop.")
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        whatsapp.disconnect()
```

---

## 2. Automatic Media Downloader

This bot automatically downloads all media (images, videos, etc.) sent to it and prints the local file path.

### Code

```python
from whatsfly import WhatsApp
import os
import time

# Define where you want to save the media
MEDIA_DIR = "./my_whatsapp_media"

def media_callback(whatsapp, event):
    # The 'MediaDownloaded' event fires after WhatsFly successfully
    # downloads and saves a media file.
    if event.get("eventType") == "MediaDownloaded":
        file_path = event.get("Path")
        info = event.get("MessageInfo", {})
        sender = info.get("PushName", "Unknown")

        print(f"Downloaded media from {sender} to: {file_path}")

if __name__ == "__main__":
    # Ensure the media directory exists
    if not os.path.exists(MEDIA_DIR):
        os.makedirs(MEDIA_DIR)

    # Initialize WhatsApp client with media_path
    # WhatsFly will automatically create subfolders like 'images', 'videos', etc.
    whatsapp = WhatsApp(media_path=MEDIA_DIR, on_event=media_callback)

    whatsapp.connect()

    print("Media downloader is running.")
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        whatsapp.disconnect()
```

### Key Points
- **`media_path`**: Passing this to the `WhatsApp` constructor enables automatic downloading of media.
- **`MediaDownloaded` event**: This event is specifically designed to notify you when a file is ready on your disk.
- **Subdirectories**: WhatsFly organizes downloads into `images/`, `audios/`, `videos/`, `documents/`, and `stickers/` within your specified `media_path`.
