package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/liangsj/vimcoplit/internal/models"
)

// Service 定义了VimCoplit的核心服务接口
type Service interface {
	// 任务相关
	CreateTask(ctx context.Context, description string) (string, error)
	GetTask(ctx context.Context, taskID string) (*Task, error)
	UpdateTask(ctx context.Context, taskID string, update TaskUpdate) error
	DeleteTask(ctx context.Context, taskID string) error

	// 文件操作
	ReadFile(ctx context.Context, path string) (string, error)
	WriteFile(ctx context.Context, path string, content string) error
	DeleteFile(ctx context.Context, path string) error

	// 命令执行
	ExecuteCommand(ctx context.Context, command string) (string, error)

	// AI交互
	GenerateResponse(ctx context.Context, taskID string, prompt string) (string, error)

	// 模型管理
	SwitchModel(ctx context.Context, modelType models.ModelType) error
	GetCurrentModel() models.ModelType
}

// Task 表示一个VimCoplit任务
type Task struct {
	ID          string
	Description string
	Status      string
	CreatedAt   int64
	UpdatedAt   int64
}

// TaskUpdate 表示任务更新
type TaskUpdate struct {
	Description *string
	Status      *string
}

// NewService 创建新的核心服务实例
func NewService() Service {
	return &serviceImpl{
		model: nil,
		mu:    &sync.RWMutex{},
	}
}

// serviceImpl 是Service接口的具体实现
type serviceImpl struct {
	model models.Model
	mu    *sync.RWMutex
}

// 实现Service接口的所有方法
func (s *serviceImpl) CreateTask(ctx context.Context, description string) (string, error) {
	// TODO: 实现创建任务的逻辑
	return "", nil
}

func (s *serviceImpl) GetTask(ctx context.Context, taskID string) (*Task, error) {
	// TODO: 实现获取任务的逻辑
	return nil, nil
}

func (s *serviceImpl) UpdateTask(ctx context.Context, taskID string, update TaskUpdate) error {
	// TODO: 实现更新任务的逻辑
	return nil
}

func (s *serviceImpl) DeleteTask(ctx context.Context, taskID string) error {
	// TODO: 实现删除任务的逻辑
	return nil
}

func (s *serviceImpl) ReadFile(ctx context.Context, path string) (string, error) {
	// TODO: 实现读取文件的逻辑
	return "", nil
}

func (s *serviceImpl) WriteFile(ctx context.Context, path string, content string) error {
	// TODO: 实现写入文件的逻辑
	return nil
}

func (s *serviceImpl) DeleteFile(ctx context.Context, path string) error {
	// TODO: 实现删除文件的逻辑
	return nil
}

func (s *serviceImpl) ExecuteCommand(ctx context.Context, command string) (string, error) {
	// TODO: 实现执行命令的逻辑
	return "", nil
}

func (s *serviceImpl) GenerateResponse(ctx context.Context, taskID string, prompt string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.model == nil {
		return "", fmt.Errorf("no model selected")
	}

	return s.model.Generate(ctx, prompt)
}

func (s *serviceImpl) SwitchModel(ctx context.Context, modelType models.ModelType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	config := models.ModelConfig{
		ModelType:   modelType,
		MaxTokens:   4096,
		Temperature: 0.7,
	}

	model, err := models.NewModel(config)
	if err != nil {
		return err
	}

	s.model = model
	return nil
}

func (s *serviceImpl) GetCurrentModel() models.ModelType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.model == nil {
		return ""
	}

	return s.model.GetModelType()
}
