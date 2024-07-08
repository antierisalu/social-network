package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	FromID   int    `json:"fromid"`
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			signal := <-SignalChan
			switch signal {
			case "privacy_updated":
				connections.Lock()
				connections.updateAllUsersStore()
				connections.Unlock()
				continue
			case "followers_updated":
				connections.Lock()
				connections.updateAllUsersStore()
				connections.Unlock()
				continue
			}
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			//log.Printf("User %v disconnected\n", connections.m[conn])
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
		switch msg.Type {
		case "login":
			connections.Lock()
			connections.m[conn] = msg.Username
			connections.Unlock()
			connections.broadcastOnlineUsers()
			connections.updateLastMsgStore(msg.Username)
			connections.updateAllUsersStore() // this is for all conns (for now..)
			// connections.updateAllUsersStore(conn) This is for a single conn
			connections.updateChatNotifStore(conn) // update the store with values from db
			log.Printf("User %s connected", msg.Username)
		case "logout":
			log.Printf("User %v logged out\n", connections.m[conn])
			connections.Lock()
			delete(connections.m, conn)
			connections.Unlock()
			connections.broadcastOnlineUsers()
			continue
/* 		case "text":
			handleTextMessage(conn, messageType, msg.Data)
		case "ping":
			handlePingMessage(conn, messageType, msg.Data) */
		case "followNotif":
			handleFollowRequest(conn, messageType, msg)
			continue
		case "followRequestNotif":
			handleFollowRequest(conn, messageType, msg)
			connections.updateAllUsersStore()
			continue
		case "acceptedFollowNotif":
			acceptedFollowRequest(conn, messageType, msg)
			connections.updateAllUsersStore()
			continue
		case "newMessageNotif":
			// handleNewMessageNotif(conn, msg.Data)
		case "groupJoinNotif":
			//handleGroupJoinNotif(conn, msg.Data)
		case "groupInviteNotif":
			//handleGroupInviteNotif(conn, msg.Data)
		case "getChatID":
			// log.Printf("User %s requested chatID", connections.m[conn])
			handleGetChatID(conn, messageType, msg.Data, msg.ID, msg.TargetID)
			continue
		case "newMessage":
			connections.handleNewMessage(conn, messageType, msg)
			// Cancel default message back to client
			continue
		case "markAsSeen":
			MarkAsSeen(msg.TargetID, msg.ID, msg.FromID)
			continue
		case "typing":
			connections.handleTyping(msg.FromID, msg.TargetID)
			continue
		default:
			log.Println("unknown message type:", msg.Type)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println("writemessage:", err)
				return
			}
		}
		// send message back to client
	}
}

func acceptedFollowRequest(conn *websocket.Conn, messageType int, msg Message) {
	fromUser, err := fetchUserByID(msg.FromID)
	if err != nil {
		fmt.Println("Error getting from email, handlefollowrequest")
		return
	}

	targetEmail, err := GetEmailFromID(msg.TargetID)
	if err != nil {
		fmt.Println("Error getting target email, handlefollowrequest")
		return
	}

	var response struct {
		Type   string `json:"type"`
		Data   string `json:"data"`
		FromID int    `json:"fromID"`
	}

	response.FromID = fromUser.ID
	response.Type = "acceptedFollowNotif"
	response.Data = fromUser.FirstName + " has accepted your request!"

	InsertNotification(fromUser.ID, response.Data, msg.Data)

	fmt.Println(response, "joptoimat")

	for usrConn, usrEmail := range connections.m {
		if targetEmail == usrEmail {
			marshaledContent, err := json.Marshal(response)
			if err != nil {
				fmt.Println("johhaidi")
			}
			// talle tahame saata
			err = usrConn.WriteMessage(messageType, marshaledContent)
			if err != nil {
				log.Println("follow notification:", err)
				// return
			}
		}
	}
}

func handleFollowRequest(conn *websocket.Conn, messageType int, msg Message) {

	// // data = link

	fromUser, err := fetchUserByID(msg.FromID)
	if err != nil {
		fmt.Println("Error getting from email, handlefollowrequest")
		return
	}

	targetEmail, err := GetEmailFromID(msg.TargetID)
	if err != nil {
		fmt.Println("Error getting target email, handlefollowrequest")
		return
	}

	var response struct {
		Type   string `json:"type"`
		Data   string `json:"data"`
		FromID int    `json:"fromID"`
	}

	followType := strings.Split(msg.Data, "_")
	if len(followType) != 2 {
		log.Printf("Invalid followtype")
		return
	}
	switch followType[0] {
	case "follow":
		response.Data = fromUser.FirstName + " has followed you!"
		response.Type = "followNotif"
	case "followRequest":
		response.Data = fromUser.FirstName + " has requested to follow you!"
		response.Type = "followRequestNotif"
	}
	response.FromID = fromUser.ID

	// fmt.Println(msg)

	InsertNotification(fromUser.ID, response.Data, msg.Data)

	for usrConn, usrEmail := range connections.m {
		if targetEmail == usrEmail {
			marshaledContent, err := json.Marshal(response)
			if err != nil {
				fmt.Println("johhaidi")
			}
			// talle tahame saata
			err = usrConn.WriteMessage(messageType, marshaledContent)
			if err != nil {
				log.Println("follow notification:", err)
				// return
			}
			return
		}
	}
}

func (c *Connections) handleTyping(fromID, targetID int) {
	targetEmail, err := GetEmailFromID(targetID)
	if err != nil {
		fmt.Println("error: couldn't get targetEmail from ID:", err)
		return
	}

	for conn, email := range c.m {
		if email == targetEmail {
			reply := struct {
				Type   string `json:"type"`
				FromID int    `json:"fromID"`
			}{
				Type:   "isTyping",
				FromID: fromID,
			}
			compiledReply, err := json.Marshal(reply)
			if err != nil {
				fmt.Println("Failed to compile array of online users to json: ", err)
			}

			err = conn.WriteMessage(1, compiledReply)
			if err != nil {
				log.Println("writemessage:", err)
			}
			break
		}
	}
}

func (c *Connections) updateChatNotifStore(ClientConn *websocket.Conn) {

	// Current userID
	clientEmail := c.m[ClientConn]
	clientID, err := GetIDFromEmail(clientEmail)
	if err != nil {
		fmt.Println("error: GetIdFromEmail: ", err)
		return
	}
	// All users
	userArr, err := FetchAllUsers(clientID)
	if err != nil {
		log.Printf("error: Failed to fetch all users: %s", err)
	}

	ChatNotifMap := make(map[int]int)

	for _, user := range userArr {
		// Check for client ID match
		if clientID == user.ID {
			continue
		}

		// Check if chat exists between users
		chatID, _ := GetChatID(clientID, user.ID)
		if chatID == -1 {
			continue
		}

		// If chat exists, check for last unseen message
		unSeenMsgID, err := GetLastUnseenMessageID(chatID)
		if err != nil {
			fmt.Println("error: getting last unseen messageID")
			continue
		}
		// No unseen messages
		if unSeenMsgID == -1 {
			continue
		}

		// Makes sure to not notify client of own messages
		messageAuthorID, err := GetMessageAuthor(unSeenMsgID)
		if err != nil {
			fmt.Println("error getting msg author: ", err)
			return
		}
		if clientID != messageAuthorID {
			ChatNotifMap[user.ID] = unSeenMsgID
		}
	}

	// Compile chatNotif map
	reply := struct {
		Type      string      `json:"type"`
		ChatNotif map[int]int `json:"chatNotif"`
	}{
		Type:      "chatNotifStore",
		ChatNotif: ChatNotifMap,
	}
	compiledReply, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("Failed to compile array of online users to json: ", err)
	}

	err = ClientConn.WriteMessage(1, compiledReply)
	if err != nil {
		log.Println("writemessage:", err)
	}
}

// Function to update allUsersStore through WS
// Pass in *ws.Conn as arg if you want to send to only that client
// Pass in no args to update allUsersStore for all connected users
// Note: currently set to all/global (bind global to registration & login for single client)
func (c *Connections) updateAllUsersStore(ClientConn ...*websocket.Conn) {
	// Compile array of []SearchData (allUsers)
	reply := struct {
		Type     string       `json:"type"`
		AllUsers []SearchData `json:"allUsers"`
	}{
		Type: "allUsers",
	}
	// Check for args (for single client)
	if len(ClientConn) > 0 {
		clientEmail := c.m[ClientConn[0]]
		clientID, err := GetIDFromEmail(clientEmail)
		if err != nil {
			fmt.Println("error: GetIdFromEmail: ", err)
			return
		}
		reply.AllUsers, err = FetchAllUsers(clientID)
		if err != nil {
			log.Printf("error: Failed to fetch all users: %s", err)
		}
		compiledReply, err := json.Marshal(reply)
		if err != nil {
			fmt.Println("Failed to compile array of online users to json: ", err)
		}
		err = ClientConn[0].WriteMessage(1, compiledReply)
		if err != nil {
			log.Println("[U1]writemessage:", err)
		}
		// fmt.Println("Send updated userList to a single client")
		return
	}
	// for all clients
	for usrConn := range c.m {
		clientEmail := c.m[usrConn]
		clientID, err := GetIDFromEmail(clientEmail)
		if err != nil {
			fmt.Println("error: GetIdFromEmail: ", err)
			return
		}
		reply.AllUsers, err = FetchAllUsers(clientID)
		if err != nil {
			log.Printf("error: Failed to fetch all users: %s", err)
		}
		compiledReply, err := json.Marshal(reply)
		if err != nil {
			fmt.Println("Failed to compile array of online users to json: ", err)
		}
		err = usrConn.WriteMessage(1, compiledReply)
		if err != nil {
			log.Println("[U2] writemessage:", err)
			fmt.Println("removing closed connection: ", clientEmail)
			delete(c.m, usrConn)

		}
	}
	// fmt.Println("Send updated userList to all users")
	// Current userID
}

func (c *Connections) updateLastMsgStore(userEmail string) {

	// Get userID from email
	userID, err := GetIDFromEmail(userEmail)
	if err != nil {
		fmt.Printf("error: failed to get ID from Email: %s", err)
	}

	// Get all last messages for ClientID & targetID
	lastMsgMap, err := GetLastMessageStore(userID)
	if err != nil {
		fmt.Printf("error: failed to get lastMessageStore for userID: %v : %s", userID, err)
	}

	// Compile map of lastMessages to json
	reply := struct {
		Type         string         `json:"type"`
		LastMsgStore map[int]string `json:"lastMsgStore"`
	}{
		Type:         "lastMsgStore",
		LastMsgStore: lastMsgMap,
	}
	compiledReply, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("Failed to compile array of online users to json: ", err)
	}

	for conn, email := range c.m {
		if email == userEmail {
			err := conn.WriteMessage(1, compiledReply)
			if err != nil {
				log.Println("writemessage:", err)
			}
			break
		}
	}

}

func (c *Connections) broadcastOnlineUsers() {
	// Get all online users
	c.Lock()
	// fmt.Println("--- Online users ---")
	onlineUserIDs := []int{}
	for _, userEmail := range c.m {
		id, err := GetIDFromEmail(userEmail)
		if err != nil {
			fmt.Println("error getting ID from email:", err)
		}
		// fmt.Println(userEmail, "ID: ", id)
		onlineUserIDs = append(onlineUserIDs, id)
	}
	// fmt.Println("--------------------")

	c.Unlock()

	// Compile array of online users to json
	reply := struct {
		Type        string `json:"type"`
		OnlineUsers []int  `json:"onlineUsers"`
	}{
		Type:        "onlineUsers",
		OnlineUsers: onlineUserIDs,
	}

	compiledReply, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("Failed to compile array of online users to json: ", err)
	}

	for usrConn := range connections.m {
		err := usrConn.WriteMessage(1, compiledReply)
		if err != nil {
			log.Println("writemessage:", err)
		}
	}

}


func (c *Connections) handleNewMessage(conn *websocket.Conn, messageType int, msg Message) {
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

			// update lastMessage list
			c.updateLastMsgStore(userEmail)

		}
		if userEmail == FromUserEmail {
			transactionFromUser = true
			// send message back to client
			err = usrConn.WriteMessage(messageType, reply)
			if err != nil {
				log.Println("writemessage:", err)
				// return
			}

			// update lastMessage list
			c.updateLastMsgStore(userEmail)

		}
	}
	// This is for handling offline notifications in the future ***TODO not implemented
	fmt.Println(transactionFromUser, transactionToUser)

}

func handleGetChatID(conn *websocket.Conn, messageType int, _ string, user1ID, user2ID int) {
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

/* func handlePingMessage(_ *websocket.Conn, messageType int, data string) {
	fmt.Println("got ping message:", messageType, data)
}

func handleTextMessage(conn *websocket.Conn, messageType int, data string) {
	fmt.Println("got text message:", messageType, data)
}
 */