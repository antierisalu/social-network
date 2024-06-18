package pkg

import "database/sql"

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// only some info for searching purposes
type SearchData struct {
	ID        int
	FirstName string
	LastName  string
	Avatar    string
}

// for registering a user
type RegisterData struct {
	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Avatar          string `json:"avatar"`
	AvatarName      string `json:"avatarName"`
	NickName        string `json:"nickName"`
	AboutMe         string `json:"aboutMe"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

// for auth checking
type User struct {
	ID          int            `json:"id"`
	FirstName   string         `json:"firstName"`
	Privacy     int            `json:"privacy"`
	LastName    string         `json:"lastName"`
	DateOfBirth sql.NullString `json:"dateOfBirth"`
	Avatar      string         `json:"avatar"`
	NickName    sql.NullString `json:"nickName"`
	AboutMe     sql.NullString `json:"aboutMe"`
	Session     sql.NullString `json:"session"`
	IsFollowing bool           `json:"isFollowing"`
	Followers   []SearchData   `json:"followers"`
	Following   []SearchData   `json:"following"`
}

type Session struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}

type PostPreview struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userID"`
	Content   string `json:"content"`
	Img       string `json:"img"`
	CreatedAt string `json:"createdAt"`
}

type Post struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userID"`
	Content   string `json:"content"`
	Img       string `json:"img"`
	CreatedAt string `json:"createdAt"`
	Privacy   int    `json:"privacy"`
	GroupID   int    `json:"groupID"`
}
