# Introduction

Welcome to the WhatsFly documentation!

WhatsFly is a powerful and easy-to-use library that enables you to interact with WhatsApp through Python. If you're familiar with Python and want to integrate WhatsApp functionalities into your projects, you've come to the right place. This library simplifies the process, allowing you to use WhatsApp with minimal effort.

## What is WhatsFly?

WhatsFly allows you to leverage the full capabilities of WhatsApp directly from your Python code. With WhatsFly, you can:

- Send and receive text messages
- Handle nd send media files (images, videos, audio)
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

| Feature                                        | Status |
|------------------------------------------------|--------|
| Multi Device                                   | ✅ |
| Send messages                                  | ✅ |
| Receive messages                               | ✅ |
| Receive media (images/audio/video/documents)   | ✅ |
| Receive location                               | ✅ |
| Send image                                     | ✅ |
| Send media (video)                             | ✅ |
| Send media (documents)                         | ✅ |
| Send media (audio)                             | ❌  |
| Send stickers                                  | ⏳  |
| Send contact cards                             | ⏳ |
| Send location                                  | ⏳ |
| Message replies                                | ⏳ |
| Join groups by invite                          | ✅ |
| Get invite for group                           | ✅ |
| Modify group name       | ✅ |
| Modify group topic | ✅ |
| Allow non-admin to edit group settings and send message (vice-versa) | ✅ |
| Add group participants                         | ⏳ |
| Kick group participants                        | ⏳ |
| Promote/demote group participants              | ⏳ |
| Mention users                                  | ⏳ |
| Mute/unmute chats                              | ⏳ |
| Block/unblock contacts                         | ⏳ |
| Get contact info                               | ⏳ |
| Get profile pictures                           | ⏳ |
| Set user status message                        | ⏳ |
| React to messages                              | ⏳ |

## Install
```bash
pip install https://github.com/Labfox/whatsfly/releases/latest/download/install.zip
```

## Usage

Here's a basic example to get you started with WhatsFly. This code demonstrates how to send a message and listen for incoming messages using WhatsFly.

### Code

```python
from whatsfly import WhatsApp
import time

def my_event_callback(event_data):
    ''' 
    Simple event callback to listen to incoming events/messages. 
    Whenever this function is called, it will retrieve the current incoming event or messages.
    '''
    print("Received event data:", event_data)

if __name__ == "__main__":

    phone = "6283139750000" # Make sure to attach country code + phone number
    message = "Hello World!"

    whatsapp = WhatsApp(event_callback=my_event_callback)

    whatsapp.connect()

    message_sent = whatsapp.sendMessage(phone=phone, message=message)
    
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

