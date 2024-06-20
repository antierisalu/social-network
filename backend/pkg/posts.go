package pkg

import (
	db "backend/pkg/db/sqlite"
	"encoding/json"
	"fmt"
	"net/http"
)

// WIP, need to figure out how to get regular posts, private posts, custom privacy posts and group posts all onto the same feed
func PostsHandler(w http.ResponseWriter, r *http.Request) {

	var groupID int
	if r.Method == "GET" {
		groupID = 0
	}
	posts, err := getPostPreviews(groupID)
	if err != nil {
		fmt.Println("PostsHandler: error ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("NewPostHandler: Autherror ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var requestBody Post
	err = decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("NewPostHandler: badRequest ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	postID, err := createPost(&requestBody, userID)
	if err != nil {
		fmt.Println("NewPostHandler: error ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(postID)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func getPostPreviews(groupID int) ([]PostPreview, error) {
	query := `SELECT id, user_id, content, media, created_at 
			  FROM posts WHERE group_id = ? 
			  ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []PostPreview
	for rows.Next() {
		var post PostPreview
		err = rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Img, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// accepts post struct pointer and returns created post ID or -1 and error
func createPost(post *Post, userID int) (int, error) {
	stmt, err := db.DB.Prepare("INSERT INTO posts (user_id, content, media, group_id, privacy) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(userID, post.Content, post.Img, post.GroupID, post.Privacy)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

/* CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    media BLOB,
    group_id INTEGER NOT NULL,
    privacy INTEGER NOT NULL CHECK (privacy BETWEEN 0 AND 2),
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
); */
