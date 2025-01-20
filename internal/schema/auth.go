package schema

type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRes struct {
	UserID string `json:"user_id"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

type GoogleOAuthReq struct {
	GoogleToken string `json:"google_token"`
}

type GoogleOAuthRes struct {
	Token string `json:"token"`
}
