package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/coltea/owlalpha/backend/internal/bootstrap"
	"github.com/coltea/owlalpha/backend/internal/model/entity"
	"github.com/coltea/owlalpha/backend/internal/service"
)

type Logic struct {
	deps *bootstrap.Dependencies
}

func New(deps *bootstrap.Dependencies) *Logic {
	return &Logic{deps: deps}
}

func (l *Logic) Login(ctx context.Context, in service.LoginInput) (*service.LoginOutput, error) {
	var user entity.User
	if err := l.deps.DB.WithContext(ctx).Where("username = ?", in.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if user.Password != in.Password {
		return nil, errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(l.deps.Config.Server.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &service.LoginOutput{
		Token: signed,
		User: &service.AuthUser{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}, nil
}

func (l *Logic) ValidateToken(ctx context.Context, tokenString string) (*service.AuthUser, error) {
	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.deps.Config.Server.JWTSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &service.AuthUser{
		ID:       uint(claims["sub"].(float64)),
		Username: claims["username"].(string),
		Role:     claims["role"].(string),
	}, nil
}
