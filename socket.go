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
func (s *socket) Register(client string, w http.ResponseWriter, r *http.Request,ReaderFunction func(msg []byte)) (err error) {
	s.connectionClient.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	con, err := s.connectionClient.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Client Connection Failed : ", err.Error())
		return
	}
	clientId := uuid.NewV4().String()
	logMsg := "New Client Connected!"
	if _, ok := s.registerClient[client]; ok {
		logMsg = "client connected !"
	}
	if s.registerClient[client].Indentifier == nil {
		var NewClient clientData
		NewClient.Indentifier = make(map[string]*websocket.Conn)
		NewClient.Indentifier[clientId] = con
		s.registerClient[client] = NewClient
	} else {
		s.registerClient[client].Indentifier[clientId] = con
	}
	go s.clientConnection(client, clientId, con,ReaderFunction)
	log.Println(logMsg)
	return
}
func (s *socket) PublishMsg(client string, msg [] byte) {
	for _, con := range s.registerClient[client].Indentifier {
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

func (s *socket) clientConnection(client, clientId string, clientCon *websocket.Conn,ReaderFunction func(asd[]byte)) {
	for {
		_, msg, err := clientCon.ReadMessage()
		if err != nil {
			log.Println(err.Error(), " : Failed to read Msg from client")
			s.disconnectClient(client, clientId)
			break
		}
		if ReaderFunction != nil {
			go ReaderFunction(msg)
		}
		if string(msg) == DISCONNECT {
			s.disconnectClient(client, clientId)
		}
	}
}

func (s *socket) disconnectClient(client, clientId string) {
	if _, ok := s.registerClient[client]; ok {
		if _, ok := s.registerClient[client].Indentifier[clientId]; ok {
			s.registerClient[client].Indentifier[clientId].Close()
			delete(s.registerClient[client].Indentifier, clientId)
			if len(s.registerClient[client].Indentifier) == 0 {
				delete(s.registerClient, client)
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
					break
					return
				}
			}
		}
		return
	}
	return errors.New("no client connected!")
}
