# Advanced Usage

WhatsFly provides access to the underlying WhatsApp protobufs, allowing you to send complex messages that aren't explicitly covered by the simple API methods.

## Using Protobufs with `sendMessage`

The `sendMessage` method's `message` parameter can accept either a string (for simple text) or a `WAWebProtobufsE2E_pb2.Message` object.

### Example: Sending a Contact Card

To send a contact card (VCard), you can construct the protobuf manually.

```python
from whatsfly.proto.waE2E import WAWebProtobufsE2E_pb2

# Create the base message object
msg = WAWebProtobufsE2E_pb2.Message()

# Add a contact message
msg.contactMessage.displayName = "John Doe"
msg.contactMessage.vcard = "BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL;type=CELL;type=VOICE;waid=1234567890:+1234567890\nEND:VCARD"

# Send it
whatsapp.sendMessage("recipient_phone_number", msg)
```

### Example: Sending a Location

```python
from whatsfly.proto.waE2E import WAWebProtobufsE2E_pb2

msg = WAWebProtobufsE2E_pb2.Message()
msg.locationMessage.degreesLatitude = -6.1754
msg.locationMessage.degreesLongitude = 106.8272
msg.locationMessage.name = "Monas"
msg.locationMessage.address = "Jakarta, Indonesia"

whatsapp.sendMessage("recipient_phone_number", msg)
```

## Custom Event Handling

You can add multiple event handlers or replace the default one.

```python
def my_custom_handler(whatsapp, event):
    if event.get("eventType") == "Connected":
        print("Connected successfully!")

whatsapp._userEventHandlers.append(my_custom_handler)
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `WHATSFLY_NO_UPDATES` | If set, WhatsFly will not attempt to check for or download binary updates from GitHub. |

## Manual Binary Management

By default, WhatsFly downloads pre-compiled Go binaries for your architecture. If you want to use your own compiled version:
1. Compile the Go backend as a shared library (`go build -buildmode=c-shared`).
2. Place the resulting file (`latest.so`, `latest.dll`, or `latest.dylib`) in the `whatsfly/dependencies/` directory.
3. Set `WHATSFLY_NO_UPDATES=1` to prevent the library from overwriting your custom binary.
