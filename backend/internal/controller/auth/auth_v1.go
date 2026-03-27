package auth

import (
	"context"

	"github.com/coltea/owlalpha/backend/api/auth/v1"
	"github.com/coltea/owlalpha/backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	data, err := service.Auth().Login(ctx, service.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.LoginRes{LoginOutput: *data}, nil
}
