# WhatsFly
[![Build](https://github.com/Labfox/whatsfly/actions/workflows/build.yml/badge.svg)](https://github.com/Labfox/whatsfly/actions/workflows/build.yml)

## Just run and have fun. Just try and go fly. 

> [!NOTE]  
> There currently isn't active development, but the project is still maintained. If you want a feature, please create an issue, I'll try to implement it as soon as possible (I usually respond within 1-2 weeks).

WhatsApp web wrapper in Python. No selenium nor gecko web driver needed. 

Setting up browser driver is tricky for python newcomers, and thus it makes your code so 'laggy' while using lots of RAM.

This project originates from [cloned-doy/whatsfly](https://github.com/cloned-doy/whatsfly)

## Documentation

https://labfox.github.io/whatsfly/

## Supported machines

The library theoretically supports every machine with Go and CGo, but if the builds fail on your machine, there are pre-built binaries auto-downloaded for the following architectures:

- Linux (amd64)
- Windows (amd64)
- macOS (amd64)
- macOS (arm64)

Additionnal architectures can be added by getting a standalone build of the go module, and adding an issue/PR to add a CI job.

## Contributing
> Please report any bug or unclear behahivor; you can also submit a PR to fix it or add a new feature.
