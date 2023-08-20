package entity

type User struct {
	ID       int
	UserID   string
	ClientID string
	RegDate  string
}

type UserResponse struct {
	TwitchId        string
	DisplayName     string
	NameColor       string
	ProfileImageUrl string
}
