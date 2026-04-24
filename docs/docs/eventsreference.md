# Events Reference

Every event not mentioned here is considered unstable and can be removed at any moment.

## Event List

| Name                 | Description                                  | Data Structure (JSON)                                                                  |
|----------------------|----------------------------------------------|-----------------------------------------------------------------------------------------|
| `Connected`          | The library is connected to WhatsApp servers | `{}`                                                                                    |
| `isLoggedIn`         | The log-in status has changed                | `{"loggedIn": bool}`                                                                    |
| `Message`            | Incoming message (Simplified)                | `{"Info": {InfoObj}, "Message": "ProtobufJSONString"}`                                  |
| `MessageJson`        | Incoming message (Full Raw Data)             | See [Message Payloads](#message-payloads)                                               |
| `MediaDownloaded`    | Fires when media is downloaded               | `{"Path": "string", "MessageInfo": {InfoObj}}`                                          |
| `Presence`           | A contact's presence has changed             | `{"from": "jid", "online": bool, "lastSeen": unix_timestamp}`                           |
| `MessageDelivered`   | A message has been delivered/read            | `{"messageIDs": ["id"], "sourceString": "jid", "timestamp": int, "type": "string"}`     |
| `AppStateSyncComplete`| The app's state has been updated            | `"name_of_sync"`                                                                        |
| `PushNameSetting`    | The account's name has been updated          | `{"timestamp": int, "action": "string", "fromFullSync": bool}`                          |
| `HistorySync`        | A part of the history has been synced        | `"filename.json"`                                                                       |
| `KeepAliveTimeout`   | Connection timeout                           | `{"errorCount": int, "lastSuccess": unix_timestamp}`                                    |
| `KeepAliveRestored`  | The library is no longer in timeout          | `{}`                                                                                    |

---

## Payload Details

### The `Info` Object
Many events include an `Info` object (of type `MessageEventCorpse` in the backend). It contains metadata about the message:

```json
{
  "ID": "ABC123DEF456",
  "MessageSource": "1234567890@s.whatsapp.net",
  "Type": "text",
  "PushName": "John Doe",
  "Timestamp": 1712345678,
  "Category": "message",
  "Multicast": false,
  "MediaType": "",
  "Ephemeral": false,
  "ViewOnce": false,
  "Edit": false,
  "Filepath": ""
}
```

### Message Payloads

#### `Message` Event
This event returns a stringified JSON of the raw WhatsApp protobuf. You need to parse it using `json.loads()`.

```python
# Example 'Message' string content for a text message:
{
  "conversation": "Hello!"
}

# Example 'Message' string content for an image message:
{
  "imageMessage": {
    "url": "https://...",
    "mimetype": "image/jpeg",
    "caption": "Check this out",
    ...
  }
}
```

#### `MediaDownloaded` Event
When `media_path` is set in the `WhatsApp` constructor, this event provides the local path to the saved file.

```json
{
  "eventType": "MediaDownloaded",
  "Path": "/absolute/path/to/whatsfly/media/images/ABC123DEF456.jpg",
  "MessageInfo": { ... }
}
```

#### `Presence` Event
```json
{
  "eventType": "Presence",
  "from": "1234567890@s.whatsapp.net",
  "online": true,
  "lastSeen": 1712345678
}
```
