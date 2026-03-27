package service

import "context"

type IHealth interface {
	Status(ctx context.Context) (*HealthStatus, error)
}

type HealthStatus struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

var localHealth IHealth

func Health() IHealth {
	return localHealth
}

func RegisterHealth(s IHealth) {
	localHealth = s
}
