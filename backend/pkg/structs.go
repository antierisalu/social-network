package pkg

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// for registering a user
type UserData struct {
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
	ID          int    `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
	Avatar      string `json:"avatar"`
	NickName    string `json:"nickName"`
	AboutMe     string `json:"aboutMe"`
	Session     string `json:"session"`
}

type Session struct {
	Token   string `json:"token"`
	Expires int    `json:"expires"`
}
