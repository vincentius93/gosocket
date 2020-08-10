package gosocket

import (
	"errors"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"sync"
)

type socket struct {
	connectionClient websocket.Upgrader
	registerClient   map[string]clientData
	lockingReWrites  *sync.RWMutex
}
type clientData struct {
	Indentifier map[string]*websocket.Conn
}
type ClientConnection = *websocket.Conn

type CallBack struct {
	currentConnection *websocket.Conn
	GoSocket          socket
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
func (s *socket) Register(Channels string, w http.ResponseWriter, r *http.Request, ReaderFunction func(msg []byte, Socket CallBack)) (clientConnection ClientConnection, err error) {
	s.connectionClient.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	con, err := s.connectionClient.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Client Connection Failed : ", err.Error())
		return nil, err
	}
	clientId := uuid.NewV4().String()
	logMsg := "New Client Connected!"
	if _, ok := s.registerClient[Channels]; ok {
		logMsg = "client connected !"
	}
	s.lockingReWrites.Lock()
	if s.registerClient[Channels].Indentifier == nil {
		var NewClient clientData
		NewClient.Indentifier = make(map[string]*websocket.Conn)
		NewClient.Indentifier[clientId] = con
		s.registerClient[Channels] = NewClient
	} else {
		s.registerClient[Channels].Indentifier[clientId] = con
	}
	s.lockingReWrites.Unlock()
	go s.clientConnection(Channels, clientId, con, ReaderFunction)
	log.Println(logMsg)
	return con, nil
}
func (s *socket) PublishMsg(Channels string, msg [] byte) {
	s.lockingReWrites.Lock()
	defer s.lockingReWrites.Unlock()
	for _, con := range s.registerClient[Channels].Indentifier {
		err := con.WriteMessage(1, msg)
		if err != nil {
			err = errors.New("Cannot send to client : " + err.Error())
		}
	}
}

func init() {
	GoSocket.registerClient = make(map[string]clientData)
	GoSocket.SetBuffer(1024, 1024)
	GoSocket.lockingReWrites = &sync.RWMutex{}
}

func (s *socket) clientConnection(Channels, clientId string, clientCon *websocket.Conn, ReaderFunction func(msg []byte, socket CallBack)) {
	for {
		_, msg, err := clientCon.ReadMessage()
		if err != nil {
			err = errors.New(err.Error() + " : Failed to read Msg from client")
			s.disconnectClient(Channels, clientId)
			break
		}
		if ReaderFunction != nil {
			newCon := CallBack{
				currentConnection: clientCon,
				GoSocket:          GoSocket,
			}
			go ReaderFunction(msg, newCon)
		}
		if string(msg) == DISCONNECT {
			s.disconnectClient(Channels, clientId)
			break
		}
	}
}

func (s *socket) disconnectClient(Channels, clientId string) {
	s.lockingReWrites.Lock()
	defer s.lockingReWrites.Unlock()
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
	s.lockingReWrites.RLock()
	defer s.lockingReWrites.RUnlock()
	if s.registerClient != nil {
		for _, value := range s.registerClient {
			for _, con := range value.Indentifier {
				err = con.WriteMessage(1, msg)
				if err != nil {
					err = errors.New("Cannot send to client : " + err.Error())
				}
			}
		}
		return
	}
	return errors.New("no channels !")
}
func (s *socket) DisconnectChannel(Channel string) (err error) {
	if _, ok := s.registerClient[Channel]; ok {
		for ClientId, _ := range s.registerClient[Channel].Indentifier {
			s.disconnectClient(Channel, ClientId)
		}
		return
	}
	err = errors.New("No Active Channel!")
	return
}

func SendToClient(clientConnection ClientConnection, msg []byte) (err error) {
	err = clientConnection.WriteMessage(1, msg)
	if err != nil {
		err = errors.New("Cannot send to client : " + err.Error())
	}
	return
}

func (c *CallBack) SendToClient(msg []byte) (err error) {
	err = c.currentConnection.WriteMessage(1, msg)
	if err != nil {
		err = errors.New("Cannot send to client : " + err.Error())
	}
	return
}
