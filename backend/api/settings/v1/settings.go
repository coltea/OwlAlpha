package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"github.com/coltea/owlalpha/backend/internal/service"
)

type GetOpenAIConfigReq struct {
	g.Meta `path:"/settings/openai" method:"get" tags:"Settings" summary:"Get OpenAI settings"`
}

type GetOpenAIConfigRes struct {
	Config *service.ModelConfig `json:"config"`
}

type CheckOpenAIConfigReq struct {
	g.Meta  `path:"/settings/openai/check" method:"post" tags:"Settings" summary:"Check OpenAI settings"`
	BaseURL string `json:"baseUrl" v:"required#请填写 Base URL"`
	APIKey  string `json:"apiKey" v:"required#请填写 API Key"`
	Model   string `json:"model" v:"required#请选择模型"`
}

type CheckOpenAIConfigRes struct {
	Message string `json:"message"`
}

type ListOpenAIModelsReq struct {
	g.Meta  `path:"/settings/openai/models" method:"post" tags:"Settings" summary:"List OpenAI models"`
	BaseURL string `json:"baseUrl" v:"required#请填写 Base URL"`
	APIKey  string `json:"apiKey" v:"required#请填写 API Key"`
}

type ListOpenAIModelsRes struct {
	Models []string `json:"models"`
}

type SaveOpenAIConfigReq struct {
	g.Meta  `path:"/settings/openai" method:"post" tags:"Settings" summary:"Save OpenAI settings"`
	BaseURL string `json:"baseUrl" v:"required#请填写 Base URL"`
	APIKey  string `json:"apiKey" v:"required#请填写 API Key"`
	Model   string `json:"model" v:"required#请选择模型"`
}

type SaveOpenAIConfigRes struct {
	Config *service.ModelConfig `json:"config"`
}
