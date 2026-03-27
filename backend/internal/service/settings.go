package service

import (
	"context"
	"time"
)

type ISettings interface {
	Get(ctx context.Context) (*ModelConfig, error)
	Check(ctx context.Context, in CheckModelConfigInput) (*CheckModelConfigOutput, error)
	ListModels(ctx context.Context, in ListModelsInput) (*ListModelsOutput, error)
	Save(ctx context.Context, in SaveModelConfigInput) (*ModelConfig, error)
}

type ModelConfig struct {
	ID        uint       `json:"id"`
	BaseURL   string     `json:"baseUrl"`
	APIKey    string     `json:"apiKey"`
	Model     string     `json:"model"`
	CheckedAt *time.Time `json:"checkedAt,omitempty"`
}

type CheckModelConfigInput struct {
	BaseURL string
	APIKey  string
	Model   string
}

type CheckModelConfigOutput struct {
	Message string `json:"message"`
}

type ListModelsInput struct {
	BaseURL string
	APIKey  string
}

type ListModelsOutput struct {
	Models []string `json:"models"`
}

type SaveModelConfigInput struct {
	BaseURL string
	APIKey  string
	Model   string
}

var localSettings ISettings

func Settings() ISettings {
	return localSettings
}

func RegisterSettings(s ISettings) {
	localSettings = s
}
