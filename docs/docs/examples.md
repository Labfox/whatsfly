# Examples

This page provides quick snippets for various tasks you can perform with WhatsFly.

## Sending Messages

### Text Message
```python
whatsapp.sendMessage("6283139750000", "Hello there!")
```

### Image with Caption
To send media, you must first upload it using `uploadFile`, then pass the returned `Upload` object to `sendMessage`.

```python
# 1. Upload the file
# kind can be: "image", "video", "audio", "document"
upload = whatsapp.uploadFile("path/to/image.jpg", kind="image")

# 2. Send the message with the upload
whatsapp.sendMessage("6283139750000", "Look at this picture!", upload=upload)
```

### Video
```python
upload = whatsapp.uploadFile("path/to/video.mp4", kind="video")
whatsapp.sendMessage("6283139750000", "Check this out", upload=upload)
```

### Document
```python
upload = whatsapp.uploadFile("path/to/report.pdf", kind="document")
whatsapp.sendMessage("6283139750000", "Here is the report", upload=upload)
```

---

## Reactions

React to a specific message using its JID and the sender's JID.

```python
# jid: the chat where the message is
# message_jid: the ID of the message to react to
# sender_jid: the JID of the person who sent the original message
# reaction: the emoji
whatsapp.sendReaction(jid, message_jid, sender_jid, "👍")
```

---

## Group Management

### Get Group Invite Link
```python
invite_link = whatsapp.getGroupInviteLink("1234567890@g.us")
print(f"Join here: {invite_link}")
```

### Join Group by Invite Link
```python
# Use the code from the end of the URL
whatsapp.joinGroupWithInviteLink("Kj9s...")
```

### Set Group Name and Topic
```python
group_jid = "1234567890@g.us"
whatsapp.setGroupName(group_jid, "New Group Name")
whatsapp.setGroupTopic(group_jid, "This is the new description")
```

### Restrict Group (Admins Only)
```python
# Only admins can send messages
whatsapp.setGroupAnnounce(group_jid, True)

# Only admins can change group info (name, icon, etc.)
whatsapp.setGroupLocked(group_jid, True)
```

### Get Group Info
```python
info = whatsapp.getGroupInfo("1234567890@g.us")
print(f"Group Name: {info['Name']}")
```

---

## Connection Management

### Check Connection Status
```python
if whatsapp.isConnected():
    print("Online")

if whatsapp.loggedIn():
    print("Logged In")
```

### Disconnect
```python
whatsapp.disconnect()
```
