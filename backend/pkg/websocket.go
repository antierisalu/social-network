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

type Connections struct {
	sync.RWMutex
	m map[*websocket.Conn]string
}

// WS
var connections = Connections{
	m: make(map[*websocket.Conn]string),
}

// var connections = struct {
// 	sync.RWMutex
// 	m map[*websocket.Conn]string
// }{m: make(map[*websocket.Conn]string)}

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
			connections.broadcastOnlineUsers()
			return
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("unmarshal:", err)
			continue
		}
		// fmt.Println("---", messageType, msg, "---")

		// fmt.Println(msg.Type)
		switch msg.Type {
		case "login":
			connections.Lock()
			connections.m[conn] = msg.Username
			connections.Unlock()
			connections.broadcastOnlineUsers()
			log.Printf("User %s connected", msg.Username)
		case "logout":
			log.Printf("User %v logged out\n", connections.m[conn])
			connections.Lock()
			delete(connections.m, conn)
			connections.Unlock()
			connections.broadcastOnlineUsers()
			continue
		case "text":
			handleTextMessage(conn, messageType, msg.Data)
		case "ping":
			handlePingMessage(conn, messageType, msg.Data)
		case "getChatID":
			// log.Printf("User %s requested chatID", connections.m[conn])
			handleGetChatID(conn, messageType, msg.Data, msg.ID, msg.TargetID)
			// Cancel default message back to client
			continue
		case "newMessage":
			handleNewMessage(conn, messageType, msg)
			// Cancel default message back to client
			continue
		case "status":

			// handleStatus(conn, messageType, msg)
			continue
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

func (c *Connections) broadcastOnlineUsers() {
	// Get all online users
	c.Lock()
	fmt.Println("all online users:")
	for _, userEmail := range c.m {
		fmt.Println("online user: ", userEmail)
	}
	c.Unlock()
}

func handleNewMessage(conn *websocket.Conn, messageType int, msg Message) {
	// ***TODO Group chats currently hardcoded for endpoint
	isGroup := false
	var pm PrivateMessage
	if err := json.Unmarshal([]byte(msg.Data), &pm); err != nil {
		log.Println("unmarshal:", err)
	}

	// Insert message to database
	createdAt, messageID, err := InsertPrivateMessage(pm.FromUserID, pm.ChatID, pm.Content, isGroup)
	if err != nil {
		fmt.Println("error Inserting private message into database!", err)
		return
	}
	// else {

	// 	// fmt.Println("inserted msg to db")
	// 	// fmt.Println("createdat: ", createdAt, " messageID: ", messageID)
	// }
	pm.Time = createdAt
	pm.MsgID = messageID
	pm.Type = "newMessage"

	reply, err := json.Marshal(pm)
	if err != nil {
		fmt.Println("ERROR")
	}
	// Check if both are online if not ***TODO add notification to offline user

	ToUserEmail, err := GetEmailFromID(pm.ToUserID)
	if err != nil {
		fmt.Println("ERROR")
	}
	// Technically dont need this can just use parent conn to reduce stack
	FromUserEmail, err := GetEmailFromID(pm.FromUserID)
	if err != nil {
		fmt.Println("ERROR")
	}

	transactionToUser := false
	transactionFromUser := false
	for usrConn, userEmail := range connections.m {
		fmt.Println("userEmail: ", userEmail)

		if userEmail == ToUserEmail {
			transactionToUser = true
			// send message back to client
			err = usrConn.WriteMessage(messageType, reply)
			if err != nil {
				log.Println("writemessage:", err)
				// return
			}
		}
		if userEmail == FromUserEmail {
			transactionFromUser = true
			// send message back to client
			err = usrConn.WriteMessage(messageType, reply)
			if err != nil {
				log.Println("writemessage:", err)
				// return
			}
		}
	}
	// This is for handling offline notifications in the future ***TODO not implemented
	fmt.Println(transactionFromUser, transactionToUser)

}

func handleGetChatID(conn *websocket.Conn, messageType int, data string, user1ID, user2ID int) {
	// fmt.Println("GO HandleGetChatID:", messageType, data)
	// fmt.Println("User IDS:", user1ID, user2ID)

	// Get Chat ID if exists
	chatID, err := GetChatID(user1ID, user2ID)
	if err != nil {
		fmt.Println("Failed to get ChatID for users", user1ID, user2ID, err)
	}

	// else {
	// 	fmt.Println("GOT THIS CHATID:", chatID)
	// }
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

	// Compile Obj structure for response
	chatIDResponse := ChatIDResponse{Type: "getChatID", ChatID: chatID}
	// Send ChatID back to client
	reply, err := json.Marshal(chatIDResponse)
	if err != nil {
		log.Println("failed to marshal reply:", err)
	}
	err = conn.WriteMessage(messageType, reply)
	if err != nil {
		log.Println("writemessage:", err)
		return
	}
	// else {
	// 	fmt.Println("Data sent to user CHATID:", chatIDResponse.ChatID)
	// }
}

func handlePingMessage(conn *websocket.Conn, messageType int, data string) {
	fmt.Println("got ping message:", messageType, data)
}

func handleTextMessage(conn *websocket.Conn, messageType int, data string) {

	fmt.Println("got text message:", messageType, data)
}
