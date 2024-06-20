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
		case "followRequestNotif":
			handleFollowRequest(msg.Data, r)
		case "newMessageNotif":
			// handleNewMessageNotif(conn, msg.Data)
		case "groupJoinNotif":
			//handleGroupJoinNotif(conn, msg.Data)
		case "groupInviteNotif":
			//handleGroupInviteNotif(conn, msg.Data)
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

func handleFollowRequest(data string, r *http.Request) {
	fromUserID, err := CheckAuth(r)
	if err != nil {
		return
	}
	fmt.Printf("fromUserID: %v\n", fromUserID)

	// data = link

	fromUser, err := fetchUserByID(fromUserID)
	if err != nil {
		fmt.Println("Error handling")
		return
	}

	content := fromUser.FirstName + " has followed you!"

	InsertNotification(fromUserID, content, data)

	// 1: Kes saatis ja kellele läheb
	// 2: Insertime DB-sse
	// 3: Teeme message ja saadame WriteMessage'iga
	// 4: Kui võetakse vastu, siis uuendame DB-s (seen DB-s) ja Frontendis
	//userID, err := CheckAuth(r)
	//connections.RLock()
	//fromUserId := connections.m[conn]
	//connections.RUnlock()

	//fmt.Printf("fromUserId: %s\n", s)

}

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

func handlePingMessage(conn *websocket.Conn, messageType int, data string) {
	fmt.Println("got ping message:", messageType, data)
}

func handleTextMessage(conn *websocket.Conn, messageType int, data string) {
	fmt.Println("got text message:", messageType, data)
}
