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
	ID       int    `json:"id"`
	TargetID int    `json:"targetid"`
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
		fmt.Println("---", messageType, msg, "---")

		fmt.Println(msg.Type)
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
		case "getChatID":
			log.Printf("User %s requested chatID", connections.m[conn])
			handleGetChatID(conn, messageType, msg.Data, msg.ID, msg.TargetID)
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

func handleGetChatID(conn *websocket.Conn, messageType int, data string, user1ID, user2ID int) {
	// fmt.Println("GO HandleGetChatID:", messageType, data)
	// fmt.Println("User IDS:", user1ID, user2ID)

	// Get Chat ID if exists
	chatID, err := GetChatID(user1ID, user2ID)
	if err != nil {
		fmt.Println("Failed to get ChatID for users", user1ID, user2ID, err)
	}
	// IF chat doesnt exist between users, create and return that chat
	if chatID == -1 {
		// fmt.Println("Chat doesn't exist between users, creating a chat for users:\n", user1ID, "\n", user2ID)
		err = InsertNewChat(user1ID, user2ID)
		if err != nil {
			fmt.Println("Failed creating a chat between these users.")
		} else {
			fmt.Println("New chat successfully created between these users.")
		}
		// Get the new chat ID
		chatID, err = GetChatID(user1ID, user2ID)
		if err != nil {
			fmt.Println("Failed to get ChatID for users", user1ID, user2ID, err)
		}
		// fmt.Println("ChatID: ", chatID)
	}

	// Send message back to client
	var reply []byte
	reply, err = json.Marshal("DATA RESPONSE HERE")
	if err != nil {
		log.Println("failed to marshal reply:", err)
	}
	err = conn.WriteMessage(messageType, reply)
	if err != nil {
		log.Println("writemessage:", err)
		return
	} else {
		fmt.Println("Data sent to user")
	}
}

func handlePingMessage(conn *websocket.Conn, messageType int, data string) {
	fmt.Println("got ping message:", messageType, data)
}

func handleTextMessage(conn *websocket.Conn, messageType int, data string) {

	fmt.Println("got text message:", messageType, data)
}
