package service

import "context"

type IAuth interface {
	Login(ctx context.Context, in LoginInput) (*LoginOutput, error)
	ValidateToken(ctx context.Context, token string) (*AuthUser, error)
}

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Token string    `json:"token"`
	User  *AuthUser `json:"user"`
}

type AuthUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var localAuth IAuth

func Auth() IAuth {
	return localAuth
}

func RegisterAuth(s IAuth) {
	localAuth = s
}
