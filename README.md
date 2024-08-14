# WhatsFly
## Just run and have fun. Just try and go fly. 

WhatsApp web wrapper in Python. No selenium nor gecko web driver needed. 

Setting up browser driver are tricky for python newcomers, and thus it makes your code so 'laggy'.

I knew that feeling. It's so painful.

So I make WhatsFly, implementing Whatsmeow --a golang based WhatsApp library. It will make his wrapper easy to use without sacrificing the speed and perfomance.

## Installation

```bash
  pip install whatsfly
```

or :
```bash
  pip install --upgrade whatsfly
```

## Usage/Examples

```javascript
from whatsfly import WhatsApp

chat = WhatsApp()

# send mesage
chat.send_message(phone="6283139750000", message="Hello World!")

# send image
chat.send_image(phone="6283139750000", image_path="path/to/image.jpg" caption="Hello World!")
```

## Features

First page on the docs

## Supported machines

| Architecture  | Status |
| ------------- | ------------- |
| Linux amd64  | ✅ |
| Linux ARM64  | ✅ |
| Linux 686  | ✅ |
| Linux 386  | ✅  |
| Windows amd64  | ✅  |
| Windows 32 bit  | soon! |
| OSX arm64  | soon! |
| OSX amd64  | soon! |

> ## Support this Project
> This project is maintained during my free time.
> If you'd like to support my work, please consider making a pull request to help fix any issues with the code.
> I would like to extend my gratitude to the open-source developers behind tls-client, tiktoken, and whatsmeow. Their work has inspired me greatly and helped me to create this project.