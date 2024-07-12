package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	db "backend/pkg/db/sqlite"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("getMessages, error getting req body:", err)
			return
		}
		var msgGet MessageGetter
		err = json.Unmarshal(body, &msgGet)
		if err != nil {
			fmt.Println("getMessages error unmarshaling,", err)
		}

		messages := GetTenMessages(msgGet.Date, msgGet.ID, msgGet.ChatID)

		jsonResponse, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			http.Error(w, "Failed to send response", http.StatusInternalServerError)
			return
		}
	}
}

func GetTenMessages(date time.Time, msgid, chatid int) []ChatMessage {
	var messages []ChatMessage

	// fmt.Printf("Getting ten messages for chatid: %v | last msgid:%v\n ", chatid, msgid)
	var query string
	if msgid == 0 { // initial load
		query = `SELECT id, content, user_id, created_at FROM chatmessages WHERE chat_id = ? and id > ? ORDER BY id DESC LIMIT 10;`
	} else { // subsequent loads
		query = `SELECT id, content, user_id, created_at FROM chatmessages WHERE chat_id = ? and id < ? ORDER BY id DESC LIMIT 10;`
	}

	rows, err := db.DB.Query(query, chatid, msgid)
	if err != nil {
		fmt.Println("GetTenMessages: error querying db: ", err)
		return []ChatMessage{}
	}

	for rows.Next() {
		var msg ChatMessage
		rows.Scan(&msg.ID, &msg.Content, &msg.User, &msg.Date)
		msg.SetUsername(db.DB)
		messages = append(messages, msg)
	}
	return messages
}

// Gets the last unseen messageID from chatmessages with chatid
// Returns -1 if no unseen messages found for the given chatID
func GetLastUnseenMessageID(chatID int) (int, error) {
	stmt := "SELECT id FROM chatmessages WHERE chat_id =? and seen = 0 ORDER BY created_at DESC LIMIT 1"
	var messageID int
	err := db.DB.QueryRow(stmt, chatID).Scan(&messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		} else {
			fmt.Println("error (GetLastUnseenMessageID): ", err)
			return -1, err
		}
	} else {
		return messageID, nil
	}

}

// Gets user_id from chatmessages with messageID
// On error returns -1, err
func GetMessageAuthor(messageID int) (int, error) {
	stmt := "SELECT user_id FROM chatmessages WHERE id = ?"
	var authorID int
	err := db.DB.QueryRow(stmt, messageID).Scan(&authorID)
	if err != nil {
		return -1, err
	}
	return authorID, nil
}

// Inserts a private message to database 'chatmessages' and returns the createdAt, message_ID, nil on success
// On error returns "ERROR", -1, err
func InsertPrivateMessage(userID, chatID int, message string, isGroup bool) (string, int, error) {
	stmt, err := db.DB.Prepare("INSERT INTO chatmessages (user_id, chat_id, content, is_group, created_at, seen) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}
	defer stmt.Close()
	//https://pkg.go.dev/database/sql#Result
	now := time.Now()
	seen := false
	result, err := stmt.Exec(userID, chatID, message, isGroup, now, seen)
	if err != nil {
		fmt.Println("Error executing statement in InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last message_ID InsertPrivateMessage:", err)
		return "ERROR", -1, err
	}
	// Update lastmessage date
	stmt, err = db.DB.Prepare(`UPDATE user_chats SET last_message = ? WHERE id = ?;`)
	if err != nil {
		fmt.Println("InsertPrivateMessage: Error Inserting LastMessage:", err)
	}
	_, err = stmt.Exec(now, chatID)
	if err != nil {
		fmt.Println("Error executing LastMessage statement in InsertPrivateMessage:", err)
	}

	return now.Format("2006-01-02 15:04:05.999999-07:00"), int(lastID), nil
}

// Creates a new database chat entry for user1 & user2
func InsertNewChat(user1, user2 int) error {
	stmt, err := db.DB.Prepare("INSERT INTO user_chats (user1, user2, last_message, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement in InsertNewChat:", err)
		return err
	}
	defer stmt.Close()
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	// Empty string for last_message
	_, err = stmt.Exec(user1, user2, "", time.Now())
	if err != nil {
		fmt.Println("Error executing statement in InsertNewChat:", err)
		return err
	}
	return nil
}

// Takes in user1 ID and user2ID to check if there is a database entry for their chat,
// IF entry exists returns the chatID if not returns -1, error
func GetChatID(userID1, userID2 int) (int, error) {
	stmt := "SELECT * FROM user_chats"
	rows, err := db.DB.Query(stmt)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			chatID      int
			user1       int
			user2       int
			createdate  string
			lastmessage string
		)
		if err := rows.Scan(&chatID, &user1, &user2, &lastmessage, &createdate); err != nil {
			fmt.Println("Error getting row values from db:", err)
			return -1, err
		}
		if (user1 == userID1 && user2 == userID2) || (user1 == userID2 && user2 == userID1) {
			//Return chatid
			return chatID, nil
		}
		if err := rows.Err(); err != nil {
			fmt.Println("Error iterating over rows:", err)
			return -1, err
		}
	}
	return -1, err
}

// Get Email From ID
func GetEmailFromID(id int) (string, error) {
	stmt := "SELECT email FROM users WHERE id = ?"
	var email string
	err := db.DB.QueryRow(stmt, id).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

// Get ID From Email
func GetIDFromEmail(email string) (int, error) {
	stmt := "SELECT id FROM users WHERE email = ?"
	var ID int
	err := db.DB.QueryRow(stmt, email).Scan(&ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

// Generate last message map (store) for clientID (PM)
func GetLastMessageStore(clientID int) (map[int]string, error) {
	lastMsgMap := make(map[int]string)

	// Fetch all rows where clientID is included
	stmt := "SELECT CASE WHEN user1 = ? THEN user2 ELSE user1 END AS other_user, last_message FROM user_chats WHERE user1 = ? OR user2 = ?"

	rows, err := db.DB.Query(stmt, clientID, clientID, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Populate the map
	for rows.Next() {
		var otherUserID int
		var lastMessage string
		if err := rows.Scan(&otherUserID, &lastMessage); err != nil {
			return nil, err
		}
		lastMsgMap[otherUserID] = lastMessage
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lastMsgMap, nil
}

// Generate last message map (store) for clientID (GM)
// Note: Includes from client msg (remove, before forwarding (ws.go))
func GetLastGroupMessageStore(clientID int) ([]int, error) {
	var groupChatIDs []int
	stmt := "SELECT g.chat_id FROM group_members gm JOIN groups g ON gm.group_id = g.id WHERE gm.user_id = ? AND (gm.chat_seen IS NULL OR gm.chat_seen = 0)"

	rows, err := db.DB.Query(stmt, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Populate the array
	for rows.Next() {
		var chatID int
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		groupChatIDs = append(groupChatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupChatIDs, nil
}

// Marks all messages seen in (user1ID + user2ID chat) before messageID(incl.)
func MarkAsSeen(messageID, user1ID, user2ID int) {
	chatID, err := GetChatID(user1ID, user2ID)
	if err != nil || chatID == -1 {
		// No chats available to perform this action
		fmt.Println("error: no chats to perform MarkAsSeen on. Ignoring..")
		return
	}

	stmt, err := db.DB.Prepare("UPDATE chatmessages SET seen = 1 WHERE chat_id = ? AND id <= ?")
	if err != nil {
		fmt.Println("Error preping DB update statement: ", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatID, messageID)
	if err != nil {
		fmt.Println("Error executing DB update statement: ", err)
	}
	// fmt.Println("marked messages as seen in chatID:", chatID, "up to message: ", messageID)
}

// Marks chat_seen to true if group_members has a row with respective groupID && userID
func MarkGroupAsSeen(groupID, userID int) error {
	var rowExists bool
	stmt := "SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)"
	err := db.DB.QueryRow(stmt, groupID, userID).Scan(&rowExists)
	if err != nil {
		return fmt.Errorf("error checking existance of row: %w", err)
	}
	if !rowExists {
		fmt.Println("no matching row found.")
		return fmt.Errorf("error no matching row found: %w", err)
	}
	// Update chat_seen to true
	// stmt2 := "UPDATE group_members SET chat_seen = 1 WHERE group_id = ? AND user_id = ?"
	// _, err = db.DB.Exec(stmt2, groupID, userID)
	// if err != nil {
	// 	return fmt.Errorf("error updating chat_seen: %w", err)
	// }

	stmt2, err := db.DB.Prepare("UPDATE group_members SET chat_seen = 1 WHERE group_id = ? AND user_id = ?")
	if err != nil {
		fmt.Println("Error preping DB update statement: ", err)
	}
	defer stmt2.Close()

	_, err = stmt2.Exec(groupID, userID)
	if err != nil {
		fmt.Println("Error executing DB update statement: ", err)
	}

	fmt.Println("chat_seen updated successfully.")
	return nil
}

// Get all members emails of a specific group with chatID
func GetGroupRecipientEmails(chatID int) ([]string, error) {
	stmt := "SELECT u.email FROM groups g JOIN group_members gm ON g.id = gm.group_id JOIN users u ON gm.user_id = u.id WHERE g.chat_id = ?"

	rows, err := db.DB.Query(stmt, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
