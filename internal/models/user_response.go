package models

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type LoginAPIResponse struct {
	Status  int           `json:"status"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    LoginResponse `json:"data"`
}

type RefreshTokenAPIResponse struct {
	Status  int                  `json:"status"`
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    RefreshTokenResponse `json:"data"`
}

type RegisterAPIResponse struct {
	Status  int              `json:"status"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    RegisterResponse `json:"data"`
}
