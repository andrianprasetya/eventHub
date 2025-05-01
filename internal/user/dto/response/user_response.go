package response

type UserResponse struct {
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	TokenType   string `json:"token_type"`
	Exp         int64  `json:"exp"`
}
