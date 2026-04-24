# Architecture

WhatsFly is designed as a bridge between Python and a high-performance WhatsApp implementation in Go.

## High-Level Overview

The project consists of two main parts:

1.  **Go Backend (`backend/`)**: Built on top of the [whatsmeow](https://github.com/tulir/whatsmeow) library. It handles the low-level WhatsApp web protocol, socket connections, encryption, and media processing.
2.  **Python Frontend (`whatsfly/`)**: Provides a user-friendly Pythonic interface. It communicates with the Go backend using Python's `ctypes` library.

```mermaid
graph LR
    UserCode[Your Python Script] --> PythonLib[WhatsFly Python Wrapper]
    PythonLib --> CTypes[ctypes / FFI]
    CTypes --> GoLib[Go Shared Library .so/.dll]
    GoLib --> WhatsMeow[whatsmeow Go Library]
    WhatsMeow --> WAServers[WhatsApp Servers]
```

## Communication (FFI)

The Go backend is compiled into a C-shared library (`.so`, `.dll`, or `.dylib`).

- **Calls from Python to Go**: Methods like `sendMessage` or `connect` call wrapped Go functions exported through the shared library.
- **Calls from Go to Python (Callbacks)**: When an event occurs (like a new message), the Go backend calls a C-compatible callback function registered by the Python wrapper.

## Threading

- **Message Thread**: When you initialize the `WhatsApp` class, a background thread is started in Python. This thread continuously calls into the Go library to process incoming network events and trigger the registered callbacks.
- **Concurrency**: The Go backend uses goroutines for internal tasks, while the Python frontend relies on standard threading to ensure the UI or main script remains responsive.

## Data Persistence

WhatsFly uses SQLite for data persistence, managed by the Go backend.
- **Session Data**: Stored in `wapp.db` within your specified `database_dir` (default is `./whatsapp/`). This keeps your login session active across restarts.
- **Media**: If enabled, media is saved directly to the filesystem in the directory specified by `media_path`.
