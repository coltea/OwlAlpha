package settings

import (
	"context"

	"github.com/coltea/owlalpha/backend/api/settings/v1"
	"github.com/coltea/owlalpha/backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) GetOpenAIConfig(ctx context.Context, _ *v1.GetOpenAIConfigReq) (res *v1.GetOpenAIConfigRes, err error) {
	config, err := service.Settings().Get(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.GetOpenAIConfigRes{Config: config}, nil
}

func (c *ControllerV1) CheckOpenAIConfig(ctx context.Context, req *v1.CheckOpenAIConfigReq) (res *v1.CheckOpenAIConfigRes, err error) {
	result, err := service.Settings().Check(ctx, service.CheckModelConfigInput{
		BaseURL: req.BaseURL,
		APIKey:  req.APIKey,
		Model:   req.Model,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CheckOpenAIConfigRes{Message: result.Message}, nil
}

func (c *ControllerV1) ListOpenAIModels(ctx context.Context, req *v1.ListOpenAIModelsReq) (res *v1.ListOpenAIModelsRes, err error) {
	result, err := service.Settings().ListModels(ctx, service.ListModelsInput{
		BaseURL: req.BaseURL,
		APIKey:  req.APIKey,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ListOpenAIModelsRes{Models: result.Models}, nil
}

func (c *ControllerV1) SaveOpenAIConfig(ctx context.Context, req *v1.SaveOpenAIConfigReq) (res *v1.SaveOpenAIConfigRes, err error) {
	config, err := service.Settings().Save(ctx, service.SaveModelConfigInput{
		BaseURL: req.BaseURL,
		APIKey:  req.APIKey,
		Model:   req.Model,
	})
	if err != nil {
		return nil, err
	}

	return &v1.SaveOpenAIConfigRes{Config: config}, nil
}
