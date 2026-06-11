package main

/*
	The frontend will understand these
	and the server transposes database records to
	the these
*/

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type UserDTO struct {
	Email    string `json:"email"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}
