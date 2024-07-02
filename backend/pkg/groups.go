package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		groups, err := getAllGroups()
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
// func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	var userID int
// 	err := decoder.Decode(&userID)
// 	if err != nil {
// 		fmt.Println("JoinGroupHandler: ", err)
// 	}

// 	err = joinGroup(userID)
// }

func createGroup(group *Group) (int, error) {
	query := `INSERT INTO groups (owner_id, name, description) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return -1, err
	}
	result, err := stmt.Exec(group.OwnerID, group.Name, group.Description)
	if err != nil {
		return -1, err
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(groupID), nil
}

func getAllGroups() ([]Group, error) {
	var groups []Group
	query := `SELECT * FROM groups`
	rows, err := db.DB.Query(query)
	if err != nil {
		return groups, err
	}

	for rows.Next() {
		var group Group
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.CreatedAt)
		if err != nil {
			fmt.Println("getAllGroups:ERROR SCANNING GROUP:", err)
			continue
		}
		groups = append(groups, group)
	}
	return groups, nil
}
