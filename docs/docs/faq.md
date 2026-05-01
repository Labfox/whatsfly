# FAQ

## Frequently Asked Questions

### 1. How do I format phone numbers?
Always use the full international format without any `+` or leading zeros.
- Correct: `6283139750000`
- Incorrect: `+6283139750000`, `083139750000`

### 2. Can I use this for bulk messaging?
While technicaly possible, WhatsApp has strict anti-spam measures. Using any third-party library for mass messaging can result in your account being banned. Always use WhatsFly responsibly and follow WhatsApp's Terms of Service.

### 3. Why does it download something on the first run?
WhatsFly needs a compiled Go backend to work. To make it easy for Python developers, we host pre-compiled binaries on GitHub and download the one matching your OS/architecture automatically.

### 4. How can I disable automatic updates?
Set the environment variable `WHATSFLY_NO_UPDATES=1` in your system or before importing the library:
```python
import os
os.environ["WHATSFLY_NO_UPDATES"] = "1"
from whatsfly import WhatsApp
```

### 5. Where is my login session stored?
By default, it is stored in a directory named `whatsapp/` in your current working directory. You can change this using the `database_dir` parameter in the `WhatsApp` constructor.

### 6. I'm getting an `OSError` related to the shared library.
This usually happens if the binary for your architecture couldn't be downloaded or is incompatible.
- Check your internet connection.
- Ensure you have the necessary permissions to write to the `whatsfly/dependencies` directory.
- If you are on an unsupported architecture, you may need to compile the Go backend manually (see [Advanced Usage](advanced.md)).

### 7. How do I stop the library?
Always call `whatsapp.disconnect()` before your script exits to ensure the Go backend closes the connection and database gracefully.

### 8. My event callback isn't being called.
- Ensure you've called `whatsapp.connect()`.
- Check if you've registered the callback either in the constructor `on_event=my_callback` or by appending to `whatsapp._userEventHandlers`.
- Make sure your script doesn't exit immediately after calling `connect()`. Use `time.sleep()` or a loop to keep it running.
