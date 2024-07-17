package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "backend/pkg/db/sqlite"
)

func NewGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, err := CheckAuth(r)
		if err != nil {
			fmt.Println("NewGroupHandler: Autherror ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var requestBody Group
		err = decoder.Decode(&requestBody)
		if err != nil {
			fmt.Println("NewGroupHandler: badRequest ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		requestBody.OwnerID = userID
		groupID, err := createGroup(&requestBody)
		if err != nil {
			fmt.Println("NewGroupHandler: error ", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		requestBody.ID = groupID

		jsonResponse, err := json.Marshal(requestBody)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
	}
}

func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("GetGroupsHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == "GET" {
		groups, err := getAllGroups(userID)
		if err != nil {
			fmt.Println("GetGroupsHandler: Error getting groups", err)
			http.Error(w, "Error getting groups", http.StatusBadRequest)
		}
		jsonResponse, err := json.Marshal(groups)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)

	}
}

// TODO: make it request to join not instajoin
func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("JoinGroupHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var ga struct {
		GroupID  int `json:"groupID"`
		Action   int `json:"action"`
		TargetID int `json:"targetID"`
	}
	err = decoder.Decode(&ga)
	if err != nil {
		fmt.Println("JoinGroupHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Println(ga)
	if ga.Action == -1 {
		err = leaveGroup(ga.TargetID, ga.GroupID)
		if err != nil {
			fmt.Println("JoinGroupHandler leave: ", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		jsonResponse, err := json.Marshal(ga)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
		return
	}
	if ga.TargetID == 0 { // set targetID to client if not not specified
		ga.TargetID = userID
	}

	err = updateGroupRelationship(ga.TargetID, ga.GroupID, ga.Action)
	if err != nil {
		fmt.Println("JoinGroupHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(ga)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("leaveGroupHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var ga struct {
		GroupID int `json:"groupID"`
	}
	err = decoder.Decode(&ga)
	if err != nil {
		fmt.Println("leaveGroupHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	err = leaveGroup(userID, ga.GroupID)
	if err != nil {
		fmt.Println("leaveGroupHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(ga)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("GetGroupHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var g struct {
		Id int `json:"groupID"`
	}
	err = decoder.Decode(&g)
	if err != nil {
		fmt.Println("GetGroupHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	group, err := GetGroup(userID, g.Id)
	if err != nil {
		fmt.Println("GetGroupHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	group.Posts, err = GetPostPreviews(group.ID, userID)
	if err != nil {
		fmt.Println("GetGroupHandler: ", err)
		http.Error(w, "Unable to get group's posts.", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(group)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

func NewEventHandler(w http.ResponseWriter, r *http.Request) {
	_, err := CheckAuth(r)
	if err != nil {
		fmt.Println("NewEventHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var event Event
	err = decoder.Decode(&event)
	if err != nil {
		fmt.Println("NewEventHandler decode error: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	eventID, err := createEvent(&event)
	if err != nil {
		fmt.Println("NewEventHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
	}

	jsonResponse, err := json.Marshal(eventID)
	if err != nil {
		fmt.Println("NewEventHandler: ", err)
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("GetEventsHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var groupID int
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&groupID)
	if err != nil {
		fmt.Println("GetEventsHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	events, err := getEvents(userID, groupID)
	if err != nil {
		fmt.Println("GetEventsHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func SendRSVPHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("SendRSVPHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var RSVP struct {
		EventID   int `json:"eventID"`
		Certainty int `json:"certainty"`
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&RSVP)
	if err != nil {
		fmt.Println("SendRSVPHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	var going int
	switch {
	case RSVP.Certainty <= 20:
		going = -1
	case RSVP.Certainty >= 80:
		going = 1
	default:
	}

	response := struct {
		EventID int `json:"eventID"`
		Going   int `json:"going"`
	}{
		EventID: RSVP.EventID,
		Going:   going,
	}

	err = sendRSVP(userID, RSVP.EventID, going)
	if err != nil {
		fmt.Println("SendRSVPHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("DeleteGroupHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var groupID struct {
		Id int `json:"groupID"`
	}
	err = decoder.Decode(&groupID)
	if err != nil {
		fmt.Println("DeleteGroupHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	err = deleteGroup(groupID.Id, userID)
	if err != nil {
		fmt.Println("DeleteGroupHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
}

func GetGroupMembersHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	_, err := CheckAuth(r)
	if err != nil {
		fmt.Println("GetGroupMembersHandler: Autherror ", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var groupID struct {
		Id int `json:"groupID"`
	}
	err = decoder.Decode(&groupID)
	if err != nil {
		fmt.Println("GetGroupMembersHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	members, err := getGroupMembers(groupID.Id)
	if err != nil {
		fmt.Println("GetGroupMembersHandler: ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
	}

	jsonResponse, err := json.Marshal(members)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
//
//	↓ HELPER FUNCTIONS ↓
//
// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

func getGroupMembers(groupID int) ([]SearchData, error) {

	var members []SearchData
	query := `SELECT u.id, firstName, lastName, avatar
			FROM group_members
			LEFT JOIN users u ON u.id = group_members.user_id
			WHERE group_id = ? AND group_members.status = 1`

	rows, err := db.DB.Query(query, groupID)
	if err != nil {
		return members, err
	}
	defer rows.Close()
	for rows.Next() {
		var member SearchData
		err = rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.Avatar)
		if err != nil {
			fmt.Println("getGroupMembers: ", err)
			continue
		}
		members = append(members, member)
	}
	return members, nil
}
func createGroup(group *Group) (int, error) {
	query := `INSERT INTO groups (owner_id, name, description, chat_id)
			VALUES (?, ?, ?, 
			COALESCE(
			(SELECT seq + 1 FROM sqlite_sequence WHERE name = 'user_chats'),
			(SELECT chat_id + 1 FROM groups WHERE id = (SELECT MAX(id) FROM groups)),
			1))`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return -1, err
	}
	result, err := stmt.Exec(group.OwnerID, group.Name, group.Description)
	if err != nil {
		return -1, err
	}

	updateSeq := `UPDATE sqlite_sequence SET seq = seq + 1 WHERE name = 'user_chats'`
	_, err = db.DB.Exec(updateSeq)
	if err != nil {
		return -1, err
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	_, err = db.DB.Exec("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, 1)", groupID, group.OwnerID)
	if err != nil {
		return -1, err
	}
	return int(groupID), nil
}

// returns all groups and whether client is joined, requested, invited, etc.
func getAllGroups(userID int) ([]Group, error) {
	var groups []Group
	query := `SELECT id, name, media, description, owner_id, groups.created_at, 
			coalesce(gm.status, -1) as joined, chat_id
			FROM groups
			LEFT JOIN group_members gm ON gm.group_id = groups.id AND gm.user_id = ?; `
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return groups, err
	}

	for rows.Next() {
		var group Group
		err = rows.Scan(&group.ID, &group.Name, &group.Media, &group.Description, &group.OwnerID, &group.CreatedAt, &group.JoinStatus, &group.ChatID)
		if err != nil {
			fmt.Println("getAllGroups:ERROR SCANNING GROUP:", err)
			continue
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// updates group relationship
// 0 = request
// 1 = join
// 2 = invite
func updateGroupRelationship(userID int, groupID int, action int) error {
	query := `INSERT OR REPLACE INTO group_members (user_id, group_id, status ) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, groupID, action)
	if err != nil {
		return err
	}
	return nil
}

// deletes group relationship (leave, unrequest, uninvite)
func leaveGroup(userID int, groupID int) error {
	query := `DELETE FROM group_members WHERE user_id = ? AND group_id = ?`
	_, err := db.DB.Exec(query, userID, groupID)
	if err != nil {
		return err
	}
	return nil
}

// userID to get followstatus of the group
func GetGroup(userID, groupID int) (Group, error) {
	query := `SELECT groups.id, name, description, media,  owner_id, groups.created_at, 
                  u.firstname || ' ' || u.lastname AS owner_name,
                 coalesce(gm.status, -1) as joined
          FROM groups
          LEFT JOIN group_members gm ON gm.group_id = groups.id AND gm.user_id = ?
          LEFT JOIN users u ON u.id = groups.owner_id
          WHERE groups.id = ?;`
	var group Group
	err := db.DB.QueryRow(query, userID, groupID).Scan(&group.ID, &group.Name, &group.Description, &group.Media, &group.OwnerID, &group.CreatedAt, &group.OwnerName, &group.JoinStatus)
	if err != nil {
		return group, err
	}

	return group, nil
}

func createEvent(event *Event) (int, error) {
	query := `INSERT INTO group_events (group_id, title, creator_id, description, date) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return -1, err
	}
	result, err := stmt.Exec(event.GroupID, event.Title, event.OwnerID, event.Description, event.Date)
	if err != nil {
		return -1, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(eventID), nil
}

// getEvents retrieves a list of events for a specific user and group.
//
// Parameters:
// - userID: The ID of the user for whom to retrieve events, this is for getting the going status.
// - groupID: The ID of the group for which to retrieve events.
//
// Returns:
// - []Event: A slice of Event objects representing the events.
// - error: An error if there was a problem retrieving the events.
func getEvents(userID, groupID int) ([]Event, error) {

	query := `SELECT 
    group_events.id, 
    group_events.title, 
    group_events.description, 
    group_events.date, 
    group_events.creator_id, 
    users.firstname || ' ' || users.lastname AS owner,
    COALESCE(gei.going, -2) AS going, 
    COUNT(CASE WHEN f1.going = 1 THEN 1 END) AS going_count,
    COUNT(CASE WHEN f1.going = 0 THEN 1 END) AS notsure_count,
    COUNT(CASE WHEN f1.going = -1 THEN 1 END) AS notgoing_count
FROM 
    group_events
LEFT JOIN 
    users ON group_events.creator_id = users.id
LEFT JOIN 
    group_event_interest gei ON group_events.id = gei.event_id AND gei.user_id = ?
LEFT JOIN 
    group_event_interest f1 ON group_events.id = f1.event_id
WHERE 
    group_events.group_id = ?
GROUP BY 
    group_events.id, 
    group_events.title, 
    group_events.description, 
    group_events.date, 
    group_events.creator_id, 
    users.firstname, 
    users.lastname, 
    gei.going
ORDER BY 
	group_events.date ASC;`

	rows, err := db.DB.Query(query, userID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event
	var oldEvents []int
	for rows.Next() {
		var event Event

		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date, &event.OwnerID, &event.OwnerName, &event.Going, &event.GoingCount, &event.NotSureCount, &event.NotGoingCount)
		if err != nil {
			fmt.Println("getEvents:ERROR SCANNING EVENT:", err)
			continue
		}

		twoHoursAfterNow := time.Now().Add(time.Hour * 2)
		eventDate, err := time.Parse(time.RFC3339, event.Date)
		if err != nil {
			continue
		}
		eventDate = eventDate.UTC()             //convert to utc
		if eventDate.Before(twoHoursAfterNow) { //if event is 2 hrs old, queue for deletion
			oldEvents = append(oldEvents, event.ID)
			continue
		}

		event.GroupID = groupID
		event.Certainty = 50
		events = append(events, event)
	}

	for _, id := range oldEvents {
		err = deleteEvent(id)
		if err != nil {
			fmt.Println("getEvents:ERROR DELETING EVENT:", err)
		}
	}
	return events, nil

}

// sendRSVP inserts or replaces a user's RSVP (going or not going) for a group event.
//
// Parameters:
// - userID: the ID of the user sending the RSVP.
// - groupID: the ID of the group.
// - going: an integer representing the user's RSVP status (1 for going, any other value for not going).
//
// Returns:
// - error: an error if there was a problem executing the SQL statement.
func sendRSVP(userID, eventID, going int) error {
	query := `INSERT OR REPLACE INTO group_event_interest (user_id, event_id, going ) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, eventID, going)
	if err != nil {
		return err
	}
	return nil
}

func deleteGroup(groupID, userID int) error {
	query := `DELETE FROM groups WHERE id = ? and owner_id = ?`
	_, err := db.DB.Exec(query, groupID, userID)
	if err != nil {
		return err
	}
	return nil
}

func deleteEvent(eventID int) error {
	query := `DELETE FROM group_events WHERE id = ?`
	_, err := db.DB.Exec(query, eventID)
	if err != nil {
		return err
	}
	return nil
}

func GetGroupOwner(groupID int) (int, string, error) {

	var ownerID int
	var ownerEmail string

	query := `SELECT id, email FROM users WHERE id = (SELECT owner_id FROM groups WHERE id = ?);`

	err := db.DB.QueryRow(query, groupID).Scan(&ownerID, &ownerEmail)

	if err != nil {
		return 0, "", err
	}

	return ownerID, ownerEmail, nil

}
