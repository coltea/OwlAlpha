package report

import (
	"context"

	"github.com/coltea/owlalpha/backend/api/report/v1"
	"github.com/coltea/owlalpha/backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) List(ctx context.Context, _ *v1.ListReq) (res *v1.ListRes, err error) {
	items, err := service.Report().List(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.ListRes{Items: items}, nil
}
