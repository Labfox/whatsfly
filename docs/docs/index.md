---
comments: true
---

# Introduction

Welcome to the WhatsFly documentation!

WhatsFly is a powerful and easy-to-use library that enables you to interact with WhatsApp through Python. If you're familiar with Python and want to integrate WhatsApp functionalities into your projects, you've come to the right place. This library simplifies the process, allowing you to use WhatsApp with minimal effort.

## What is WhatsFly?

WhatsFly allows you to leverage the full capabilities of WhatsApp directly from your Python code. With WhatsFly, you can:

- Send and receive text messages
- Handle and send media files (images, videos, audio)
- Receive notifications
- And much more

WhatsFly provides a Pythonic interface, making it easy to incorporate WhatsApp functionalities into your Python applications without dealing with the complexities of lower-level implementations.

## Why Use WhatsFly?

WhatsFly offers a streamlined and efficient way to integrate WhatsApp into your Python projects. By avoiding the use of a WebDriver, WhatsFly operates faster and more resource-efficiently. This means:

- **Improved Performance:** Directly interacting with WhatsApp's underlying protocols ensures quicker response times compared to the overhead of WebDriver-based solutions.
- **Resource Optimization:** By not relying on a WebDriver, WhatsFly consumes fewer system resources, making it suitable for both small-scale applications and large-scale deployments.
- **Reliability:** Minimizing dependencies on external tools reduces the chances of encountering issues related to browser updates or compatibility.

## Current Features

✅: Works
❌: Broke
⏳: Soon
🔧: Can work with some tinkering

|                               Feature                                |                          Status                          |
|:--------------------------------------------------------------------:|:--------------------------------------------------------:|
|                             Multi Device                             |                            ✅                             |
|                            Send messages                             |                            ✅                             |
|                           Receive messages                           |                            ✅                             |
|             Receive media (images/audio/video/documents)             |                            ✅                             |
|                           Receive location                           |                            ✅                             |
|                              Send image                              |                            ✅                             |
|                          Send media (video)                          |                            ✅                             |
|                        Send media (documents)                        |                            ✅                             |
|                          Send media (audio)                          |                            ✅                             |
|                            Send stickers                             |                  ⏳: update the uploader                  |
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
|                                Polls                                 | ⏳: create function + vote funtion, correctly reads votes |
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
If go is found in the path, the binaries will be built dynamically
```bash
pip install whatsfly-Labfox
```

## Usage

Here's a basic example to get you started with WhatsFly. This code demonstrates how to send a message and listen for incoming messages using WhatsFly.

### Code

```python
from whatsfly import WhatsApp
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

    whatsapp = WhatsApp(on_event=my_event_callback)

    whatsapp.connect()

    message_sent = whatsapp.sendMessage(phone, message, False)
    
    time.sleep(5 * 60)  # Listen for messages for 5 minutes

    whatsapp.disconnect()
```

Warning: The first time you will start the library, it will compile or download binaries for your machine, so it could take long depending on your internet connexion.

### Explanation

1. **Event Callback Function:**
   - `my_event_callback(event_data)` handles incoming events and simply prints the event data to the console.

2. **Main Program Flow:**
   - The phone number (with country code) and the message to be sent are defined.
   - An instance of `WhatsApp` is created with the event callback function.
   - The script connects to WhatsApp using the `connect()` method. At this point it should show a QR code, scan it with your phone (on Connected Devices)
   - A message is sent to the specified phone number using `sendMessage(phone, message)`.
   - The script listens for incoming messages for 5 minutes using `time.sleep(5 * 60)`.
   - Finally, it disconnects from WhatsApp using the `disconnect()` method.

