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
	FromID   int    `json:"fromid"`
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
		case "followRequestNotif":
			handleFollowRequest(conn, messageType, msg)
			break
		case "newMessageNotif":
			// handleNewMessageNotif(conn, msg.Data)
		case "groupJoinNotif":
			//handleGroupJoinNotif(conn, msg.Data)
		case "groupInviteNotif":
			//handleGroupInviteNotif(conn, msg.Data)
		case "getChatID":
			// log.Printf("User %s requested chatID", connections.m[conn])
			handleGetChatID(conn, messageType, msg.Data, msg.ID, msg.TargetID)
			// Cancel default message back to client
			continue
		case "newMessage":
			handleNewMessage(conn, messageType, msg)
			// Cancel default message back to client
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

func handleFollowRequest(conn *websocket.Conn, messageType int, msg Message) {

	// // data = link

	fromUser, err := fetchUserByID(msg.FromID)
	if err != nil {
		fmt.Println("Error getting from email, handlefollowrequest")
		return
	}
	fmt.Println(fromUser)

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

	response.Data = fromUser.FirstName + " has followed you!"
	response.FromID = fromUser.ID
	response.Type = "followRequestNotif"

	fmt.Println(msg)

	InsertNotification(fromUser.ID, response.Data, msg.Data)

	for usrConn, usrEmail := range connections.m {
		fmt.Println("usrEmail: ", usrEmail)
		fmt.Println("targetEmail: ", targetEmail)
		fmt.Println("fromEmail: ", fromUser.Email)
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

// 1: Kes saatis ja kellele läheb
// 2: Insertime DB-sse
// 3: Teeme message ja saadame WriteMessage'iga
// 4: Kui võetakse vastu, siis uuendame DB-s (seen DB-s) ja Frontendis
//userID, err := CheckAuth(r)
//connections.RLock()
//fromUserId := connections.m[conn]
//connections.RUnlock()

//fmt.Printf("fromUserId: %s\n", s)

/*
{
type:"groupinv",
fromUserId:69,
toUserId:68,
groupId: 420 / nil,
}

for

*/
//notif link: follow_{user_id} > follow_4
//group invite example: groupinvite_{groupid}_{userid}

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
