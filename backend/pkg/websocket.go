package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections = struct {
	sync.RWMutex
	m map[*websocket.Conn]string
}{m: make(map[*websocket.Conn]string)}

type Message struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Username string `json:"username"`
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %v disconnected\n", connections.m[conn])
			connections.Lock()
			delete(connections.m, conn)
			connections.Unlock()
			return
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("unmarshal:", err)
			continue
		}
		
		switch msg.Type {
		case "login":
			connections.Lock()
			connections.m[conn] = msg.Username
			connections.Unlock()
			log.Printf("User %s connected", msg.Username)
		case "text":
			handleTextMessage(conn, messageType, msg.Data)
		case "ping":
			handlePingMessage(conn, messageType, msg.Data)
		default:
			log.Println("unknown message type:", msg.Type)
		}

		// send message back to client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("writemessage:", err)
			return
		}
	}
}

func handlePingMessage(conn *websocket.Conn, messageType int, data string) {
	fmt.Println("got ping message:", messageType, data)
}

func handleTextMessage(conn *websocket.Conn, messageType int, data string) {
	

	fmt.Println("got text message:", messageType, data)
}
