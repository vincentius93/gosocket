# GoSocket

GoSocket is package for go webSocket protocol. This package provide a 
client connection management.

### Resources
```
https://github.com/gorilla/websocket
https://github.com/satori/go.uuid
```
### Installation
```
go get github.com/vincentius93/gosocke
```

### Features

- Client allows to disconnect server connection
- Manage client connection 
- Send Broadcast to all connected client
- Publish Message to all connected client on channels

## Usage
#### User connection
```
GoSocket.Register("Channel1",httpResponseWriter,httpRequest,nil)
```

#### User connection with received message function handler
```
GoSocket.Register("Channel1",httpResponseWriter,httpRequest, func(msg []byte) {
            // do some action here
            fmt.Println(string(msg))
	})
```

#### Publish Message
This function allow server to publish message to all connected client on the same channel
````
GoSocket.PublishMsg("Channel1",[]byte("hello world"))
````

#### Broadcast Message
This function allow server to broadcast message to all connected client on all channels
````
GoSocket.Broadcast([]byte("Hello Worlds"))
````

#### Client disconnect server connection
GoSocket allow client to disconnect server connection by simply send this command from 
client side
````
socket.send("diconnect")
````