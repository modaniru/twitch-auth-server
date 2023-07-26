package client

type ValidateResponse struct {
	ClientId  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserId    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
}
type UsersInformation struct {
	Data []UserInformation `json:"data"`
}

type UserInformation struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflienImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

type UsersColor struct {
	Data []*UserColor `json:"data"`
}

type UserColor struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
	Color     string `json:"color"`
}

type AppCredentials struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
