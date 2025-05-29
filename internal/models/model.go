package models

import (
	"context"
	"fmt"
)

// ModelType 定义支持的模型类型
type ModelType string

const (
	ModelTypeClaude   ModelType = "claude-3-sonnet-20240229"
	ModelTypeDoubao   ModelType = "doubao"
	ModelTypeDeepSeek ModelType = "deepseek"
)

// Model 定义了AI模型的接口
type Model interface {
	// Generate 生成响应
	Generate(ctx context.Context, prompt string) (string, error)

	// GetModelType 返回模型类型
	GetModelType() ModelType
}

// ModelConfig 定义了模型配置
type ModelConfig struct {
	APIKey      string
	ModelType   ModelType
	MaxTokens   int
	Temperature float64
}

// NewModel 创建新的模型实例
func NewModel(config ModelConfig) (Model, error) {
	switch config.ModelType {
	case ModelTypeClaude:
		return newClaudeModel(config)
	case ModelTypeDoubao:
		return newDoubaoModel(config)
	case ModelTypeDeepSeek:
		return newDeepSeekModel(config)
	default:
		return nil, fmt.Errorf("unsupported model type: %s", config.ModelType)
	}
}

// claudeModel Claude模型实现
type claudeModel struct {
	config ModelConfig
}

func newClaudeModel(config ModelConfig) (Model, error) {
	return &claudeModel{
		config: config,
	}, nil
}

func (m *claudeModel) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: 实现Claude API调用
	return "", nil
}

func (m *claudeModel) GetModelType() ModelType {
	return m.config.ModelType
}

// doubaoModel 豆包模型实现
type doubaoModel struct {
	config ModelConfig
}

func newDoubaoModel(config ModelConfig) (Model, error) {
	return &doubaoModel{
		config: config,
	}, nil
}

func (m *doubaoModel) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: 实现豆包API调用
	return "", nil
}

func (m *doubaoModel) GetModelType() ModelType {
	return m.config.ModelType
}

// deepSeekModel DeepSeek模型实现
type deepSeekModel struct {
	config ModelConfig
}

func newDeepSeekModel(config ModelConfig) (Model, error) {
	return &deepSeekModel{
		config: config,
	}, nil
}

func (m *deepSeekModel) Generate(ctx context.Context, prompt string) (string, error) {
	// TODO: 实现DeepSeek API调用
	return "", nil
}

func (m *deepSeekModel) GetModelType() ModelType {
	return m.config.ModelType
}
