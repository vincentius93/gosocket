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
go get github.com/vincentius93/gosocket
```

### Features

- Client allow to disconnect server connection
- Manage client connection 
- Send Broadcast to all connected client
- Publish Message to all connected client on channels
- Send Message to specific client
- Server Allow to disconnect all client from specific channel

## Usage
#### User connection
```
_,err := gosocket.GoSocket.Register("Channel1",w,r, nil)
if err != nil {
    // do some error handling here
}
```

#### User connection with received message function handler
This package also support to callback a function when receive new message from client
```
ClientConnection,err := gosocket.GoSocket.Register("Channel1",w,r, func(msg []byte, socket gosocket.CallBack) {
            //do some action here 
	})
	if err != nil {
		//do some error handling here
	}
```

#### Publish Message
This function allow server to publish message to all connected client on the same channel
````
GoSocket.PublishMsg("Channel1",[]byte("hello world"))
````

#### Broadcast Message
This function allow server to broadcast message to all connected client on all channels
````
GoSocket.Broadcast([]byte("Hello World"))
````

#### Client disconnect server connection
GoSocket allow client to disconnect server connection by simply send this command from 
client side
````
socket.send("diconnect")
````

#### Server disconnect all client connection from channel
GoSocket allow server to disconnect all client connection from channel with this command
````
err := gosocket.DisconnectChannel(channel)
if err != nil {
    // do some error handling here
}
````

#### Send to specific client
This package also support for server to send message to specific client
````
err := gosocket.SendToClient(ClientConnection,[]byte(val))
if err != nil {
    // do some error handling here
}
````