package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

type Connections struct {
	sync.RWMutex
	m  map[*websocket.Conn]string
	rm map[string]*websocket.Conn
}

// WS
var connections = Connections{
	m:  make(map[*websocket.Conn]string),
	rm: make(map[string]*websocket.Conn),
}

// var connections = struct {
// 	sync.RWMutex
// 	m map[*websocket.Conn]string
// }{m: make(map[*websocket.Conn]string)}

type Message struct {
	Type           string `json:"type"`
	Data           string `json:"data"`
	Username       string `json:"username"`
	ID             int    `json:"id"`
	TargetID       int    `json:"targetid"`
	FromID         int    `json:"fromid"`
	NotificationID int    `json:"notificationid"`
	IsGroup        bool   `json:"isgroup"`
	GroupID        int    `json:"groupid"`
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
			delete(connections.rm, connections.m[conn])
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
			connections.rm[msg.Username] = conn
			connections.Unlock()
			connections.broadcastOnlineUsers()
			connections.updateLastMsgStore(msg.Username)
			connections.updateAllUsersStore() // this is for all conns (for now..)
			// connections.updateAllUsersStore(conn) This is for a single conn
			connections.updateChatNotifStore(conn)      // update the PM store with values from db
			connections.updateGroupChatNotifStore(conn) // update the GM store with values from db
			//log.Printf("User %s connected", msg.Username)
		case "logout":
			log.Printf("User %v logged out\n", connections.m[conn])
			connections.Lock()
			delete(connections.m, conn)
			delete(connections.rm, connections.m[conn])
			connections.Unlock()
			connections.broadcastOnlineUsers()
			continue
		case "followNotif":
			handleFollowRequest(conn, messageType, msg)
			continue
		case "followRequest":
			handleFollowRequest(conn, messageType, msg)
			connections.updateAllUsersStore()
			continue
		case "acceptedFollow":
			acceptedFollowRequest(conn, messageType, msg)
			connections.updateAllUsersStore()
			continue
		case "declinedRequest":
			declinedRequest(conn, messageType, msg)
			continue
		case "acceptedGroupRequest":
			acceptedGroupRequest(conn, messageType, msg)
			connections.updateAllUsersStore()
			continue
		case "acceptedGroupInvite":
			acceptedGroupInvite(conn, messageType, msg)
			connections.updateAllUsersStore()
			continue
		case "cancelRequest":
			cancelFollowRequest(conn, messageType, msg)
			continue
		case "clearNotif":
			clearNotification(msg.FromID)
		case "clearSingleNotif":
			clearSingleNotifForSelf(conn, messageType, msg)

		case "groupRequest":
			handleGroupRequest(conn, messageType, msg)
		case "groupInvite":
			handleGroupInvite(conn, messageType, msg)

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
		case "markGroupAsSeen":
			// Group Chat seen states in group_members (chat_seen)
			err := MarkGroupAsSeen(msg.TargetID, msg.FromID)
			if err != nil {
				log.Println(err)
			}
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
		fmt.Println("Err getting from email, handlefollowrequest")
		return
	}

	targetEmail, err := GetEmailFromID(msg.TargetID)
	if err != nil {
		fmt.Println("Error getting target email, handlefollowrequest")
		return
	}

	var response struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      bool   `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
	}

	response.FromID = fromUser.ID
	response.Type = "acceptedFollow"
	response.Content = fromUser.FirstName + " has accepted your request!"

	notification, err := InsertNotification(msg.TargetID, response.Content, msg.Data)
	if err != nil {
		fmt.Println("Error inserting notification, acceptedFollowRequest")
	}

	response.ID = notification.ID
	response.Link = notification.Link
	response.Seen = notification.Seen
	response.CreatedAt = notification.CreatedAt

	err = clearSingleNotification(msg.NotificationID)
	if err != nil {
		log.Printf("Error clearing notification acceptedfollowrequest")
	}

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


func clearSingleNotifForSelf(conn *websocket.Conn, messageType int, msg Message) {
	notificID, err := strconv.Atoi(msg.Data)
	if err != nil {
		log.Printf("Failed to convert msg data to int clear single notif - websocket.go")
	}
	err = clearSingleNotification(notificID)
	if err != nil {
		log.Printf("Failed to clear single notif from db - websocket.go")
	}

}

func cancelFollowRequest(conn *websocket.Conn, messageType int, msg Message) {
	targetEmail, err := GetEmailFromID(msg.TargetID)
	if err != nil {
		fmt.Println("Error getting target email, cancelFollowRequest")
		return
	}

	var response struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      string `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
	}

	response.Type = msg.Type

	fromIDStr := strconv.Itoa(msg.FromID)
	followRequestLink := "followRequest_" + fromIDStr

	// find notification id based on userid and link
	followRequestNotification, err := GetNotificationBasedOnLink(followRequestLink)
	if err != nil {
		fmt.Println("Failed getting notification based on link - cancelFollowRequest")
	}

	err = clearSingleNotification(followRequestNotification.ID)
	if err != nil {
		log.Printf("Failed to clear single notif from db cancelFollowRequest - websocket.go")
	}

	response.ID = followRequestNotification.ID

	for usrConn, usrEmail := range connections.m {
		if targetEmail == usrEmail {
			marshaledContent, err := json.Marshal(response)
			if err != nil {
				fmt.Println("johhaidi")
			}
			// talle tahame saata
			err = usrConn.WriteMessage(messageType, marshaledContent)
			if err != nil {
				log.Println("remove notification websocket:", err)
				// return
			}
		}
	}
}

func declinedRequest(conn *websocket.Conn, messageType int, msg Message) {

	err := clearSingleNotification(msg.NotificationID)
	if err != nil {
		log.Printf("Failed to clear single notif from db cancelFollowRequest - websocket.go")
	}
}

func handleFollowRequest(conn *websocket.Conn, messageType int, msg Message) {

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
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      bool   `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
	}

	exists, err := CheckDuplicateNotification(msg.TargetID, msg.Data)
	if err != nil {
		fmt.Println("Error checking duplicate notification, handlefollowrequest")

	}
	if exists {
		fmt.Println("User already has notification, handle follow request")
		return
	}

	followType := strings.Split(msg.Data, "_")
	if len(followType) != 2 {
		log.Printf("Invalid followtype")
		return
	}
	switch followType[0] {
	case "follow":
		response.Content = fromUser.FirstName + " has followed you!"
		response.Type = "follow"
	case "followRequest":
		response.Content = fromUser.FirstName + " has requested to follow you!"
		response.Type = "followRequest"
	}
	response.FromID = fromUser.ID

	notification, err := InsertNotification(msg.TargetID, response.Content, msg.Data)
	if err != nil {
		fmt.Println("Error inserting notification, handlefollowrequest")
	}

	response.ID = notification.ID
	response.Link = notification.Link
	response.Seen = notification.Seen
	response.CreatedAt = notification.CreatedAt

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

// ChatNotifStore For Group messages
func (c *Connections) updateGroupChatNotifStore(ClientConn *websocket.Conn) {
	// Current userID
	clientEmail := c.m[ClientConn]
	clientID, err := GetIDFromEmail(clientEmail)
	if err != nil {
		fmt.Println("error: GetIdFromEmail: ", err)
		return
	}

	// Retrieve chat_id from the groups table where
	// the corresponding user_id in the group_members table matches the function argument and chat_seen is NULL or 0
	notifChatIDs, err := GetLastGroupMessageStore(clientID)
	if err != nil {
		fmt.Println("error: GetLastGroupMessageStore: ", err)
		return
	}

	// All the group Chat IDs that have unresolved notification (group chat notifications)
	var groupChatIDs []int

	// Check if the last message is from client
	for _, chatid := range notifChatIDs {
		// Note: for groupchats seen is always 0 so it just gets the last messageID from that groupchat
		messageID, err := GetLastUnseenMessageID(chatid)
		if err != nil {
			fmt.Println("error: getting last unseen messageID")
			continue
		}

		// Get Message author ID Makes sure to not notify client of own messages
		messageAuthorID, err := GetMessageAuthor(messageID)
		if err != nil {
			fmt.Println("error getting msg author: ", err)
			return
		}

		if clientID != messageAuthorID {
			groupChatIDs = append(groupChatIDs, chatid)
		}

	}
	/* 	fmt.Println("GROUPCHATIDS WITH UNRESOLVED NOTIF:")
	   	fmt.Println(groupChatIDs) */
	// Compile chatNotif response for groupchats
	reply := struct {
		Type      string `json:"type"`
		ChatNotif []int  `json:"chatNotif"`
	}{
		Type:      "groupChatNotifStore",
		ChatNotif: groupChatIDs,
	}

	compiledReply, err := json.Marshal(reply)
	if err != nil {
		fmt.Println("Failed to compile array of unseen group messages to json: ", err)
	}

	err = ClientConn.WriteMessage(1, compiledReply)
	if err != nil {
		log.Println("writemessage:", err)
	}

}

// ChatNotifStore For Private messages
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
		fmt.Println("Failed to compile array of chatNotifStore to json: ", err)
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
	isGroup := false
	isGroup = msg.IsGroup
	var pm PrivateMessage
	var gm GroupMessage
	var groupRecipients []string
	var reply []byte
	var ToUserEmail, FromUserEmail string
	fmt.Println("MESSAGE FROM GROUP?:", isGroup)
	if isGroup {
		if err := json.Unmarshal([]byte(msg.Data), &gm); err != nil {
			log.Println("unmarshal:", err)
		}

		// Insert message to database
		createdAt, messageID, err := InsertPrivateMessage(gm.FromUserID, gm.ChatID, gm.Content, isGroup)
		if err != nil {
			fmt.Println("error Inserting private message into database!", err)
			return
		}

		gm.Time = createdAt
		gm.MsgID = messageID
		gm.Type = "newGroupMessage"

		reply, err = json.Marshal(gm)
		if err != nil {
			fmt.Println("ERROR")
		}

		// Group Messages (TO BE OPTIMIZED)
		groupRecipients, err = GetGroupRecipientEmails(gm.ChatID)
		//setUnseen(groupRecipients, groupID) ***Lukas
		if err != nil {
			fmt.Printf("error: GetGroupRecipientEmails(%s): %s", groupRecipients, err)
			return
		} else {
			fmt.Println("Group members:", groupRecipients)
		}
	} else {
		if err := json.Unmarshal([]byte(msg.Data), &pm); err != nil {
			log.Println("unmarshal:", err)
		}

		// Insert message to database
		createdAt, messageID, err := InsertPrivateMessage(pm.FromUserID, pm.ChatID, pm.Content, isGroup)
		if err != nil {
			fmt.Println("error Inserting private message into database!", err)
			return
		}

		pm.Time = createdAt
		pm.MsgID = messageID
		pm.Type = "newMessage"

		reply, err = json.Marshal(pm)
		if err != nil {
			fmt.Println("ERROR")
		}

		ToUserEmail, err = GetEmailFromID(pm.ToUserID)
		if err != nil {
			fmt.Println("ERROR GetEmailFromID(ToUserID)")
		}
		// Technically dont need this can just use parent conn to reduce stack
		FromUserEmail, err = GetEmailFromID(pm.FromUserID)
		if err != nil {
			fmt.Println("ERROR GetEmailFromID(FromUserID)")
		}
	}
	fmt.Println(pm, gm)
	// else {
	// STILL NEED TO HANDLE SETTING CHAT_SEEN VALUES BACK TO 0 or NULL
	// WHEN TO DO IT?
	// PROBABLY ON NEW MESSAGE?
	// 	// fmt.Println("inserted msg to db")
	// 	// fmt.Println("createdat: ", createdAt, " messageID: ", messageID)
	// }
	// pm.Time = createdAt
	// pm.MsgID = messageID
	// pm.Type = "newMessage"

	// reply, err := json.Marshal(pm)
	// if err != nil {
	// 	fmt.Println("ERROR")
	// }

	// ToUserEmail, err := GetEmailFromID(pm.ToUserID)
	// if err != nil {
	// 	fmt.Println("ERROR")
	// }
	// // Technically dont need this can just use parent conn to reduce stack
	// FromUserEmail, err := GetEmailFromID(pm.FromUserID)
	// if err != nil {
	// 	fmt.Println("ERROR")
	// }
	for usrConn, userEmail := range connections.m {
		// Group Messages
		fmt.Println("ISGROUP:", isGroup)
		if isGroup {
			for _, targetEmail := range groupRecipients {
				fmt.Println("USEREMAIL:", userEmail, "TARGET EMAIL:", targetEmail)
				if userEmail == targetEmail {
					//update seen ID for each recipient in database
					err := usrConn.WriteMessage(messageType, reply)
					if err != nil {
						log.Println("writemessage:", err)
					}
					fmt.Println("SENDING GROUPCHAT DATA OUT TO USER:", userEmail)
					// update lastMessage list
					c.updateLastMsgStore(userEmail)
					continue
				}
			}

		} else {
			// Private Messages
			if userEmail == ToUserEmail || userEmail == FromUserEmail {
				// send message back to client
				err := usrConn.WriteMessage(messageType, reply)
				if err != nil {
					log.Println("writemessage:", err)
				}

				// update lastMessage list
				c.updateLastMsgStore(userEmail)
			}
		}

	}
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

func handleGroupRequest(conn *websocket.Conn, messageType int, msg Message) {

	ownerID, ownerEmail, err := GetGroupOwner(msg.GroupID)
	if err != nil {
		fmt.Println("Failed to get group owner")
		return
	}

	exists, err := CheckDuplicateNotification(ownerID, msg.Data) //check if such a notification already exists
	if err != nil {
		fmt.Println("Error checking duplicate notification, handlefollowrequest")

	}
	if exists {
		fmt.Println("User already has notification, handle follow request")
		return
	}

	toConn := connections.rm[ownerEmail]
	if toConn == nil {
		fmt.Println("Failed to get connection for user", ownerEmail)
		return
	}

	response := struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      bool   `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
		GroupID   int    `json:"groupID"`
	}{
		Type:    "groupRequest",
		GroupID: msg.GroupID,
		FromID:  msg.FromID,
	}

	fromUser, err := fetchUserByID(msg.FromID) //get from users info for notification content
	if err != nil {
		fmt.Println("Error getting from email, handlefollowrequest")
		return
	}

	linkArr := strings.Split(msg.Data, "_")
	if len(linkArr) != 3 {
		log.Printf("Invalid followtype")
		return
	}

	group, err := GetGroup(msg.TargetID, msg.GroupID)

	if err != nil {
		fmt.Println("Error getting group, handlefollowrequest")
		return
	}

	switch linkArr[0] {
	case "groupRequest":
		response.Content = fromUser.FirstName + " has requested to join your group " + group.Name + "!"
		response.Type = "groupRequest"
	}

	notification, err := InsertNotification(ownerID, response.Content, msg.Data)
	if err != nil {
		fmt.Println("Error inserting notification, handlefollowrequest")
	}

	response.ID = notification.ID
	response.Link = notification.Link
	response.Seen = notification.Seen
	response.CreatedAt = notification.CreatedAt

	reply, err := json.Marshal(response)

	if err != nil {
		log.Println("failed to marshal reply:", err)
		return
	}

	toConn.WriteMessage(messageType, reply)
}

/*  */

func handleGroupInvite(conn *websocket.Conn, messageType int, msg Message) {

	exists, err := CheckDuplicateNotification(msg.TargetID, msg.Data) //check if such a notification already exists
	if err != nil {
		fmt.Println("Error checking duplicate notification, handlefollowrequest")

	}
	if exists {
		fmt.Println("User already has notification, handle follow request")
		return
	}

	targetEmail, err := GetEmailFromID(msg.TargetID) //get target email to send notification
	if err != nil {
		fmt.Println("Error getting target email, handleGroupInvite")
		return
	}

	toConn := connections.rm[targetEmail] //get targets websocket connection
	if toConn == nil {
		fmt.Println("Failed to get connection for user", targetEmail)
		return
	}

	response := struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      bool   `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
		GroupID   int    `json:"groupID"`
	}{
		Type:    "groupInvite",
		GroupID: msg.GroupID,
		FromID:  msg.FromID,
	}

	fromUser, err := fetchUserByID(msg.FromID) //get from users info for notification content
	if err != nil {
		fmt.Println("Error getting from email, handlefollowrequest")
		return
	}

	linkArr := strings.Split(msg.Data, "_")
	if len(linkArr) != 3 {
		log.Printf("Invalid followtype")
		return
	}

	group, err := GetGroup(msg.TargetID, msg.GroupID)

	if err != nil {
		fmt.Println("Error getting group, handlefollowrequest")
		return
	}

	switch linkArr[0] {
	case "groupInvite":
		response.Content = fromUser.FirstName + " has Invited you to join the group " + group.Name + "!"
		response.Type = "groupInvite"
	}

	notification, err := InsertNotification(msg.TargetID, response.Content, msg.Data)
	if err != nil {
		fmt.Println("Error inserting notification, handlefollowrequest")
	}

	response.ID = notification.ID
	response.Link = notification.Link
	response.Seen = notification.Seen
	response.CreatedAt = notification.CreatedAt

	reply, err := json.Marshal(response)

	if err != nil {
		log.Println("failed to marshal reply:", err)
		return
	}

	toConn.WriteMessage(messageType, reply)
}

func acceptedGroupRequest(conn *websocket.Conn, messageType int, msg Message) {
	fromUser, err := fetchUserByID(msg.FromID)
	if err != nil {
		fmt.Println("Err getting from email, acceptedGroupRequest")
		return
	}

	targetEmail, err := GetEmailFromID(msg.TargetID)
	if err != nil {
		fmt.Println("Error getting target email, acceptedGroupRequest")
		return
	}

	group, err := GetGroup(msg.TargetID, msg.GroupID)

	if err != nil {
		fmt.Println("Error getting group, handlefollowrequest")
		return
	}

	toConn := connections.rm[targetEmail]

	var response struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      bool   `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
	}

	response.FromID = fromUser.ID
	response.Type = "acceptedGroupRequest"
	response.Content = fromUser.FirstName + " has accepted your request to join " + group.Name + "!"

	notification, err := InsertNotification(msg.TargetID, response.Content, msg.Data)
	if err != nil {
		fmt.Println("Error inserting notification, acceptedGroupRequest")
	}

	response.ID = notification.ID
	response.Link = notification.Link
	response.Seen = notification.Seen
	response.CreatedAt = notification.CreatedAt

	err = clearSingleNotification(msg.NotificationID)
	if err != nil {
		log.Printf("Error clearing notification acceptedGroupRequest")
	}
	reply, err := json.Marshal(response)

	if err != nil {
		log.Println("failed to marshal reply:", err)
		return
	}

	toConn.WriteMessage(messageType, reply)


}
func acceptedGroupInvite(conn *websocket.Conn, messageType int, msg Message) {
	fmt.Print("acceptedGroupInvite")
	fromUser, err := fetchUserByID(msg.FromID)
	if err != nil {
		fmt.Println("Err getting from email, acceptedGroupRequest")
		return
	}

	targetEmail, err := GetEmailFromID(msg.TargetID)
	if err != nil {
		fmt.Println("Error getting target email, acceptedGroupRequest")
		return
	}

	group, err := GetGroup(msg.TargetID, msg.GroupID)

	if err != nil {
		fmt.Println("Error getting group, handlefollowrequest")
		return
	}

	toConn := connections.rm[targetEmail]

	var response struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Link      string `json:"link"`
		Seen      bool   `json:"seen"`
		CreatedAt string `json:"createdAt"`
		FromID    int    `json:"fromID"`
	}

	response.FromID = fromUser.ID
	response.Type = "acceptedGroupRequest"
	response.Content = fromUser.FirstName + " has accepted your invite to join " + group.Name + "!"

	notification, err := InsertNotification(msg.TargetID, response.Content, msg.Data)
	if err != nil {
		fmt.Println("Error inserting notification, acceptedGroupRequest")
	}

	response.ID = notification.ID
	response.Link = notification.Link
	response.Seen = notification.Seen
	response.CreatedAt = notification.CreatedAt

	err = clearSingleNotification(msg.NotificationID)
	if err != nil {
		log.Printf("Error clearing notification acceptedGroupRequest")
	}
	reply, err := json.Marshal(response)

	if err != nil {
		log.Println("failed to marshal reply:", err)
		return
	}

	toConn.WriteMessage(messageType, reply)
}
