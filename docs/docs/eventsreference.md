---
comments: true
---

Every event not mentioned here is considered unstable and can be removed at any moment.


| Name                 | Description                                  | Members                                                                     |
|----------------------|----------------------------------------------|-----------------------------------------------------------------------------|
| AppStateSyncComplete | The app's state has been updated             | name (```str```)                                                            |
| Connected            | The library is connected to whatsapp servers | None                                                                        |
| PushNameSetting      | The account's name has been updated          | timestamp (unix timestamp), action (```str```), fromFullSync (```bool```)   |
| Message              | Incoming message                             | info (```dict```), message(```dict```)                                      |
| MessageRead          | A message has been read                      | messageID (```list```), sourceString(```str```), timestamp (unix timestamp) |
| Presence             | A contact's presence has changed             | from (```str```), online (```bool```), lastSeen (unix timestamp, optional)  |
| HistorySync          | A part of the history has been synced        | filename (```str```)                                                        |
| KeepAliveTimeout     | Connection timeout                           | errorCount (```int```), lastSuccess (unix timestamp)                        |
| KeepAliveRestored    | The library is not longer in timeout         | None                                                                        |
| isLoggedIn           | The log-in status has changed                | loggedIn (```bool```)                                                       |
| MediaDownloaded      | Fires when a media is downloaded             | path(```str```), associatedMessageInfo(```dict```)                          |