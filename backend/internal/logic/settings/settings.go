package settings

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/coltea/owlalpha/backend/internal/bootstrap"
	"github.com/coltea/owlalpha/backend/internal/model/entity"
	"github.com/coltea/owlalpha/backend/internal/service"
)

const requestTimeout = 20 * time.Second

type Logic struct {
	deps   *bootstrap.Dependencies
	client *http.Client
}

type openAIModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type openAIErrorResponse struct {
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func New(deps *bootstrap.Dependencies) *Logic {
	return &Logic{
		deps: deps,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (l *Logic) Get(ctx context.Context) (*service.ModelConfig, error) {
	var config entity.ModelConfig
	if err := l.deps.DB.WithContext(ctx).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &service.ModelConfig{
				BaseURL: l.deps.Config.OpenAI.BaseURL,
				APIKey:  l.deps.Config.OpenAI.APIKey,
				Model:   l.deps.Config.OpenAI.Model,
			}, nil
		}
		return nil, err
	}

	return toServiceModelConfig(config), nil
}

func (l *Logic) Check(ctx context.Context, in service.CheckModelConfigInput) (*service.CheckModelConfigOutput, error) {
	normalized, err := normalizeAndValidate(in.BaseURL, in.APIKey, in.Model)
	if err != nil {
		return nil, err
	}

	if err := l.requestChatCompletion(ctx, normalized.BaseURL, normalized.APIKey, normalized.Model); err != nil {
		return nil, err
	}

	return &service.CheckModelConfigOutput{Message: "配置检查通过，可以保存。"}, nil
}

func (l *Logic) ListModels(ctx context.Context, in service.ListModelsInput) (*service.ListModelsOutput, error) {
	baseURL, apiKey, err := normalizeBaseCredentials(in.BaseURL, in.APIKey)
	if err != nil {
		return nil, err
	}

	models, err := l.requestModels(ctx, baseURL, apiKey)
	if err != nil {
		return nil, err
	}

	return &service.ListModelsOutput{Models: models}, nil
}

func (l *Logic) Save(ctx context.Context, in service.SaveModelConfigInput) (*service.ModelConfig, error) {
	normalized, err := normalizeAndValidate(in.BaseURL, in.APIKey, in.Model)
	if err != nil {
		return nil, err
	}

	if err := l.requestChatCompletion(ctx, normalized.BaseURL, normalized.APIKey, normalized.Model); err != nil {
		return nil, err
	}

	now := time.Now()
	db := l.deps.DB.WithContext(ctx)

	var config entity.ModelConfig
	err = db.First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config = entity.ModelConfig{
			BaseURL:   normalized.BaseURL,
			APIKey:    normalized.APIKey,
			Model:     normalized.Model,
			CheckedAt: &now,
		}
		if err := db.Create(&config).Error; err != nil {
			return nil, err
		}
		return toServiceModelConfig(config), nil
	}
	if err != nil {
		return nil, err
	}

	config.BaseURL = normalized.BaseURL
	config.APIKey = normalized.APIKey
	config.Model = normalized.Model
	config.CheckedAt = &now
	if err := db.Save(&config).Error; err != nil {
		return nil, err
	}

	return toServiceModelConfig(config), nil
}

func (l *Logic) requestModels(ctx context.Context, baseURL, apiKey string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+"/models", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求模型列表失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, buildRemoteError("拉取模型列表失败", resp.StatusCode, body)
	}

	var payload openAIModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("解析模型列表失败: %w", err)
	}

	set := make(map[string]struct{}, len(payload.Data))
	for _, item := range payload.Data {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			continue
		}
		set[id] = struct{}{}
	}

	models := make([]string, 0, len(set))
	for id := range set {
		models = append(models, id)
	}
	sort.Strings(models)

	if len(models) == 0 {
		return nil, errors.New("模型列表为空，请确认接口权限与返回格式")
	}

	return models, nil
}

func (l *Logic) requestChatCompletion(ctx context.Context, baseURL, apiKey, model string) error {
	payload := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "ping",
			},
		},
		"max_tokens":  1,
		"temperature": 0,
		"stream":      false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求模型检查失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return buildRemoteError("模型检查失败", resp.StatusCode, respBody)
	}

	return nil
}

func normalizeAndValidate(baseURL, apiKey, model string) (*service.SaveModelConfigInput, error) {
	normalizedBaseURL, normalizedAPIKey, err := normalizeBaseCredentials(baseURL, apiKey)
	if err != nil {
		return nil, err
	}

	normalizedModel := strings.TrimSpace(model)
	if normalizedModel == "" {
		return nil, errors.New("请选择模型")
	}

	return &service.SaveModelConfigInput{
		BaseURL: normalizedBaseURL,
		APIKey:  normalizedAPIKey,
		Model:   normalizedModel,
	}, nil
}

func normalizeBaseCredentials(baseURL, apiKey string) (string, string, error) {
	normalizedBaseURL := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if normalizedBaseURL == "" {
		return "", "", errors.New("请填写 Base URL")
	}
	if !strings.HasPrefix(normalizedBaseURL, "http://") && !strings.HasPrefix(normalizedBaseURL, "https://") {
		return "", "", errors.New("Base URL 必须是有效的 http 或 https 地址")
	}

	normalizedAPIKey := strings.TrimSpace(apiKey)
	if normalizedAPIKey == "" {
		return "", "", errors.New("请填写 API Key")
	}

	return normalizedBaseURL, normalizedAPIKey, nil
}

func buildRemoteError(prefix string, statusCode int, body []byte) error {
	var payload openAIErrorResponse
	if err := json.Unmarshal(body, &payload); err == nil && payload.Error != nil && strings.TrimSpace(payload.Error.Message) != "" {
		return fmt.Errorf("%s: %s", prefix, strings.TrimSpace(payload.Error.Message))
	}

	snippet := strings.TrimSpace(string(body))
	if snippet == "" {
		return fmt.Errorf("%s: 远端返回状态码 %d", prefix, statusCode)
	}
	if len(snippet) > 240 {
		snippet = snippet[:240]
	}

	return fmt.Errorf("%s: 远端返回状态码 %d, %s", prefix, statusCode, snippet)
}

func toServiceModelConfig(config entity.ModelConfig) *service.ModelConfig {
	return &service.ModelConfig{
		ID:        config.ID,
		BaseURL:   config.BaseURL,
		APIKey:    config.APIKey,
		Model:     config.Model,
		CheckedAt: config.CheckedAt,
	}
}
