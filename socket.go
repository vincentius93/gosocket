package gosocket

import (
	"errors"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

type socket struct {
	connectionClient websocket.Upgrader
	registerClient   map[string]clientData
}
type clientData struct {
	Indentifier map[string]*websocket.Conn
}
type ClientConnection = *websocket.Conn

const (
	DISCONNECT = "disconnect"
)

var GoSocket socket

func (s *socket) SetBuffer(read, write int) {
	s.connectionClient.ReadBufferSize = read
	s.connectionClient.WriteBufferSize = write
	if s.registerClient == nil {
		s.registerClient = make(map[string]clientData)
	}
}
func (s *socket) Register(Channels string, w http.ResponseWriter, r *http.Request,ReaderFunction func(msg []byte,ClientCon ClientConnection)) (clientConnection ClientConnection,err error) {
	s.connectionClient.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	con, err := s.connectionClient.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Client Connection Failed : ", err.Error())
		return nil,err
	}
	clientId := uuid.NewV4().String()
	logMsg := "New Client Connected!"
	if _, ok := s.registerClient[Channels]; ok {
		logMsg = "client connected !"
	}
	if s.registerClient[Channels].Indentifier == nil {
		var NewClient clientData
		NewClient.Indentifier = make(map[string]*websocket.Conn)
		NewClient.Indentifier[clientId] = con
		s.registerClient[Channels] = NewClient
	} else {
		s.registerClient[Channels].Indentifier[clientId] = con
	}

	go s.clientConnection(Channels, clientId, con,ReaderFunction)
	log.Println(logMsg)
	return con,nil
}
func (s *socket) PublishMsg(Channels string, msg [] byte) {
	for _, con := range s.registerClient[Channels].Indentifier {
		err := con.WriteMessage(1, msg)
		if err != nil {
			log.Println("Cannot send to client : ", err.Error())
		}
	}
	log.Println("Message sent!")
}

func init() {
	GoSocket.registerClient = make(map[string]clientData)
	GoSocket.SetBuffer(1024, 1024)
}

func (s *socket) clientConnection(Channels, clientId string, clientCon *websocket.Conn,ReaderFunction func(msg []byte,clientConnection ClientConnection)) {
	for {
		_, msg, err := clientCon.ReadMessage()
		if err != nil {
			log.Println(err.Error(), " : Failed to read Msg from client")
			s.disconnectClient(Channels, clientId)
			break
		}
		if ReaderFunction != nil {
			go ReaderFunction(msg,clientCon)
		}
		if string(msg) == DISCONNECT {
			s.disconnectClient(Channels, clientId)
		}
	}
}

func (s *socket) disconnectClient(Channels, clientId string) {
	if _, ok := s.registerClient[Channels]; ok {
		if _, ok := s.registerClient[Channels].Indentifier[clientId]; ok {
			s.registerClient[Channels].Indentifier[clientId].Close()
			delete(s.registerClient[Channels].Indentifier, clientId)
			if len(s.registerClient[Channels].Indentifier) == 0 {
				delete(s.registerClient, Channels)
			}
			log.Println("User disconnected!")
		}
	}
}
func (s *socket) Broadcast(msg []byte) (err error) {
	if s.registerClient != nil {
		for _, value := range s.registerClient {
			for _, con := range value.Indentifier {
				err = con.WriteMessage(1, msg)
				if err != nil {
					log.Println("Cannot send to client : ", err.Error())
				}
			}
		}
		return
	}
	return errors.New("no channels !")
}

func SendToClient(clientConnection ClientConnection,msg []byte)(err error){
	err = clientConnection.WriteMessage(1, msg)
	if err != nil {
		log.Println("Cannot send to client : ", err.Error())
	}
	return
}