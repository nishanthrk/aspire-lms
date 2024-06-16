package dto

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"  validate:"required"`
	Platform   string `json:"-" validate:"required"`
}

type TokenDetail struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	UserType     string `json:"user_type"`
	TokenExpires int64  `json:"token_expires"`
}

type LoginResponse struct {
	Data   TokenDetail `json:"data"`
	Status int         `json:"status"`
}
