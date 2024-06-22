package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	db "backend/pkg/db/sqlite"
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
	/*
		avatarHandler := Avatars{
			NewAvatar: AvatarImgData{
				Base64String: requestBody.Img,
				FileName:     userData.AvatarName,
			},
			Dir: "avatars",
		} */

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

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckAuth(r)
	if err != nil {
		fmt.Println("NewCommentHandler: Autherror ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestBody Comment
	err = decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("NewCommentHandler: badRequest ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	commentID, err := createComment(&requestBody, userID)
	if err != nil {
		fmt.Println("NewCommentHandler: error ", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(commentID)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

// func getPostPreviews(groupID int) ([]PostPreview, error) {
// 	query := `SELECT id, user_id, content, media, created_at
// 			  FROM posts WHERE group_id = ?
// 			  ORDER BY created_at DESC`

// 	rows, err := db.DB.Query(query, groupID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var posts []PostPreview
// 	for rows.Next() {
// 		var post PostPreview
// 		err = rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Img, &post.CreatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		posts = append(posts, post)
// 	}
// 	return posts, nil
// }

func getPostPreviews(groupID int) ([]PostPreview, error) {
	postsQuery := `SELECT id, user_id, content, media, created_at
                   FROM posts
                   WHERE group_id = ?
                   ORDER BY created_at DESC`

	commentsQuery := `SELECT c.id, c.user_id, c.post_id, c.content, c.created_at,
                            u.FirstName, u.LastName, u.Avatar
                      FROM comments c
                      JOIN users u ON c.user_id = u.id
                      WHERE c.post_id IN (SELECT id FROM posts WHERE group_id = ?)
                      ORDER BY c.created_at ASC`

	// Fetch posts
	postRows, err := db.DB.Query(postsQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer postRows.Close()

	var posts []PostPreview
	for postRows.Next() {
		var post PostPreview
		var img sql.NullString
		err = postRows.Scan(&post.ID, &post.UserID, &post.Content, &img, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		post.Img = img.String
		post.Comments = []Comment{} // Initialize the comments slice
		posts = append(posts, post)
	}

	// Fetch comments
	commentRows, err := db.DB.Query(commentsQuery, groupID)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	commentsMap := make(map[int][]Comment)
	for commentRows.Next() {
		var comment Comment
		err = commentRows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt,
			&comment.User.FirstName, &comment.User.LastName, &comment.User.Avatar)
		if err != nil {
			return nil, err
		}
		commentsMap[comment.PostID] = append(commentsMap[comment.PostID], comment)
	}

	// Merge comments into their respective posts
	for i := range posts {
		if comments, found := commentsMap[posts[i].ID]; found {
			posts[i].Comments = comments
		}
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

func createComment(comment *Comment, userID int) (int, error) {
	stmt, err := db.DB.Prepare("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(userID, comment.PostID, comment.Content)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// CREATE TABLE IF NOT EXISTS comments (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     user_id INTEGER NOT NULL,
//     post_id INTEGER NOT NULL,
//     content TEXT NOT NULL,
//     created_at DATE NOT NULL DEFAULT CURRENT_DATE,
//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
//     FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
// );

/* CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    media BLOB,
    group_id INTEGER NOT NULL,
    privacy INTEGER NOT NULL CHECK (privacy BETWEEN 0 AND 2),
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
); */
