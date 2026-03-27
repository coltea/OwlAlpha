package health

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/coltea/owlalpha/backend/api/health/v1"
	"github.com/coltea/owlalpha/backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) Status(ctx context.Context, _ *v1.StatusReq) (res *v1.StatusRes, err error) {
	data, err := service.Health().Status(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.StatusRes{HealthStatus: *data}, nil
}


var _ g.Meta
