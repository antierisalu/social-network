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
		GroupID int `json:"groupID"`
		Action  int `json:"action"`
	}
	err = decoder.Decode(&ga)
	if err != nil {
		fmt.Println("JoinGroupHandler: ", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	err = updateGroupRelationship(userID, ga.GroupID, ga.Action)
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

	group, err := getGroup(userID, g.Id)
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

// returns all groups and whether client is joined, requested, invited, etc.
func getAllGroups(userID int) ([]Group, error) {
	var groups []Group
	query := `SELECT id, name, description, owner_id, groups.created_at, 
			coalesce(gm.status, -1) as joined
			FROM groups
			LEFT JOIN group_members gm ON gm.group_id = groups.id AND gm.user_id = ?; `
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return groups, err
	}

	for rows.Next() {
		var group Group
		err = rows.Scan(&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.CreatedAt, &group.JoinStatus)
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

func getGroup(userID, groupID int) (Group, error) {
	query := `SELECT groups.id, name, description, owner_id, groups.created_at, 
                  u.firstname || ' ' || u.lastname AS owner_name,
                 coalesce(gm.status, -1) as joined
          FROM groups
          LEFT JOIN group_members gm ON gm.group_id = groups.id AND gm.user_id = ?
          LEFT JOIN users u ON u.id = groups.owner_id
          WHERE groups.id = ?;`
	var group Group
	err := db.DB.QueryRow(query, userID, groupID).Scan(&group.ID, &group.Name, &group.Description, &group.OwnerID, &group.CreatedAt, &group.OwnerName, &group.JoinStatus)
	if err != nil {
		return group, err
	}

	return group, nil
}
