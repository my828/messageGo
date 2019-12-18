package handlers

import (
	"assignments-my828/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.
// A simple store to store all the connections

// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

type SocketStore struct {
	Connections map[int64]*websocket.Conn
	lock        sync.Mutex
	Context     *Context
}

// !!!!!!!!!!!!
// type Notifier struct {
// 	user map[int64][]*socketStore
// }

// Thread-safe method for inserting a connection
func (s *SocketStore) InsertConnection(conn *websocket.Conn, id int64) {
	s.lock.Lock()
	// insert socket connection
	s.Connections[id] = conn
	// s.Connections = append(s.Connections, conn)
	s.lock.Unlock()
}

func (s *SocketStore) RemoveConnection(id int64) {
	s.lock.Lock()
	delete(s.Connections, id)
	s.lock.Unlock()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// This function's purpose is to reject websocket upgrade requests if the
		// origin of the websockete handshake request is coming from unknown domains.
		// This prevents some random domain from opening up a socket with your server.
		// TODO: make sure you modify this for your HW to check if r.Origin is your host
		return true
	},
}

type rabbitMes struct {
	UserIDs []int64 `json:"userIDs, omitempty"`
}

func (s *SocketStore) ConsumeMessage(msgs <-chan amqp.Delivery) error {
	for msg := range msgs {
		s.lock.Lock()
		userIDs := &rabbitMes{}
		err := json.Unmarshal(msg.Body, userIDs)
		if err != nil {
			return err
		}
		if len(userIDs.UserIDs) == 0 {
			for id, conn := range s.Connections {
				err = conn.WriteMessage(TextMessage, msg.Body)
				if err != nil {
					s.RemoveConnection(id)
					return err
				}
			}
		} else {
			for _, id := range userIDs.UserIDs {
				conn, errBool := s.Connections[id]
				if !errBool {
					s.RemoveConnection(id)
				}
				err = conn.WriteMessage(TextMessage, msg.Body)
				if err != nil {
					s.RemoveConnection(id)
					conn.Close()
				}

			}
		}
		s.lock.Unlock()
	}
	return nil
}

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

func (s *SocketStore) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {

	state := &SessionState{}
	_, err := sessions.GetState(r, s.Context.Key, s.Context.SessionStore, state)
	if err != nil {
		fmt.Printf("Error getting session state/session unauthorized %v", http.StatusUnauthorized)
		return
	}
	// handle the websocket handshake
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}
	// Insert our connection onto our datastructure for ongoing usage
	s.InsertConnection(conn, state.User.ID)

	go (func(conn *websocket.Conn, userID int64, ws *SocketStore) {
		defer conn.Close()
		defer s.RemoveConnection(state.User.ID)

		for {
			messageType, _, err := conn.ReadMessage()
			if messageType == CloseMessage {
				fmt.Println("Close message received.")
				break
			} else if err != nil {
				fmt.Println("Error reading message.")
				break
			}
		}
	})(conn, state.User.ID, s)
}
