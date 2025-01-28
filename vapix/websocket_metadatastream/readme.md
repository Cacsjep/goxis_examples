# VAPIX WebSocket Metadata Stream Example

This example demonstrates how to use the `vapix` package to connect to an Axis device and subscribe to its metadata stream using WebSockets. 
The application leverages the `acapapp` package for logging and cleanup and the `vapix` package for configuring and consuming metadata events.

[AXIS Metadata Via Websocket](https://help.axis.com/en-us/axis-os-knowledge-base#metadata-via-websocket)

## Requirements

- Axis device with VAPIX WebSocket metadata support (AXIS OS 10.11 or newer)
- The following Go packages:
  - `github.com/Cacsjep/goxis/pkg/acapapp`
  - `github.com/Cacsjep/goxis/pkg/vapix`

## Installation

Clone the example repo or include the necessary packages in your Go program:

```shell
go get github.com/Cacsjep/goxis/pkg/acapapp
go get github.com/Cacsjep/goxis/pkg/vapix
```

### Build example with goxisbuilder
``` shell
go install github.com/Cacsjep/goxisbuilder@latest
git clone https://github.com/Cacsjep/goxis_examples
cd goxis_examples
goxisbuilder.exe -appdir "./vapix/websocket_metadatastream"
```