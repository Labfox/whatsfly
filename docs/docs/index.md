---
comments: true
---

# Introduction

Welcome to the WhatsFly documentation!

This is mostly a work in progress, if you don't understand something, or some beahivor isn't described correctly, please open an issue.

## What is WhatsFly?

WhatsFly allows you to interface with Whatsapp from Python. It is pretty easy to use, and usually works with the latest version. If you have a more advanced use case, please use [whatsmeow](https://github.com/tulir/whatsmeow) (Go) or [baileys](https://github.com/WhiskeySockets/Baileys) (Typescript)


## Why Use WhatsFly?

- You want to send a text-based notification from your python script
- You want to trigger something from whatsapp
- You want to download all media sent to you on whatsapp
- You want to create a simple whatsapp bot

## Current Features

✅: Works
❌: Broken
⏳: Easily doable with some tinkering (create an issue if you want the feature)
🔧: Needs tinkering

|                               Feature                                |                           Status                          |
|:--------------------------------------------------------------------:|:---------------------------------------------------------:|
|                             Multi Device                             |                            ✅                             |
|                            Send messages                             |                            ✅                             |
|                           Receive messages                           |                            ✅                             |
|             Receive media (images/audio/video/documents)             |                            ✅                             |
|                           Receive location                           |                            ✅                             |
|                              Send image                              |                            ✅                             |
|                          Send media (video)                          |                            ✅                             |
|                        Send media (documents)                        |                            ✅                             |
|                          Send media (audio)                          |                            ✅                             |
|                            Send stickers                             |                  ⏳                  |
|                          Send contact cards                          |                            ✅                             |
|                            Send location                             |                            ✅                             |
|                           Message replies                            |                            ✅                             |
|                        Join groups by invite                         |                            ✅                             |
|                         Get invite for group                         |                            ✅                             |
|                          Modify group name                           |                            ✅                             |
|                          Modify group topic                          |                            ✅                             |
| Allow non-admin to edit group settings and send message (vice-versa) |                            ✅                             |
|                            Get Group info                            |                            ✅                             |
|                        Add group participants                        |                            ⏳                             |
|                       Kick group participants                        |                            ⏳                             |
|                  Promote/demote group participants                   |                            ⏳                             |
|                            Mention users                             |                            ✅                             |
|                          Mute/unmute chats                           |                            ⏳                             |
|                        Block/unblock contacts                        |                            ⏳                             |
|                         Get profile pictures                         |                            ⏳                             |
|                       Set user status message                        |                            ⏳                             |
|                             Create Group                             |                            ⏳                             |
|                          Create Newsletter                           |                            ⏳                             |
|                                Polls                                 | ⏳ |
|                          Receive Reactions                           |                            ✅                             |
|                           React to message                           |                            ✅                             |
|                          Follow newsletter                           |                            ⏳                             |
|                         Get business profile                         |                            ⏳                             |
|                         Get contact QR link                          |                            ⏳                             |
|                 Get group info from invite and link                  |                            ⏳                             |
|                    Get group participants request                    |                            ⏳                             |
|                          Get joined Groups                           |                            ⏳                             |
|                      Get community participants                      |                            ⏳                             |
|                         Get newsletter info                          |                            ⏳                             |
|                         Get privacy settings                         |                            ⏳                             |
|                       Get profile picture info                       |                            ⏳                             |
|                            Set/get status                            |                            ⏳                             |
|                       Get groups of community                        |                            ⏳                             |
|                            Get user info                             |                            ⏳                             |
|                      Get if user is on whatsapp                      |                            ⏳                             |
|                        Join group with invite                        |                            ✅                             |
|                             Leave group                              |                            ⏳                             |
|                      Link group with community                       |                            ⏳                             |
|                             Mark as read                             |                            ⏳                             |
|                            Send presence                             |                            ⏳                             |
|                                                                      |                                                          |


## Install

```bash
pip install whatsfly-Labfox
```

## Usage

Here's a basic example to get you started with WhatsFly. This code demonstrates how to send a message and listen for incoming messages using WhatsFly.

### Code

```python
from whatsfly import WhatsApp # May take some time or throw an exception
import time
import pprint

def my_event_callback(whatsapp, event_data):
    ''' 
    Simple event callback to listen to incoming events/messages. 
    Whenever this function is called, it will retrieve the current incoming event or messages.
    '''
    pprint.pprint("Received event data:", event_data)

if __name__ == "__main__":

    phone = "6283139750000" # Make sure to attach country code + phone number
    message = "Hello World!"

    whatsapp = WhatsApp(on_event=my_event_callback) # The client object

    whatsapp.connect() # Will print a QR Code to the terminal

    message_sent = whatsapp.sendMessage(phone, message, False)
    
    time.sleep(5 * 60)  # Listen for messages for 5 minutes

    whatsapp.disconnect()
```

Warning: The first time you will start the library, and once every month, Whatsfly will download the binaries from Github. You can disable this beahivor by setting the `WHATSFLY_NO_UPDATES` environment variable.
