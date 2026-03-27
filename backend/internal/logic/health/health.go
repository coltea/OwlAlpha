package health

import (
	"context"

	"github.com/coltea/owlalpha/backend/internal/service"
)

type Logic struct{}

func New() *Logic {
	return &Logic{}
}

func (l *Logic) Status(ctx context.Context) (*service.HealthStatus, error) {
	return &service.HealthStatus{
		Status:   "ok",
		Database: "connected",
		Redis:    "connected",
	}, nil
}
