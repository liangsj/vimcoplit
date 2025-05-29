package mcp

import (
	"context"
	"fmt"
	"time"
)

// Tool 表示一个 MCP 工具
type Tool struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Parameters  []ToolParameter   `json:"parameters"`
	ServerID    string            `json:"server_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// ToolParameter 表示工具参数
type ToolParameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     any    `json:"default,omitempty"`
}

// Server 表示一个 MCP 服务器
type Server struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	URL         string            `json:"url"`
	Type        ServerType        `json:"type"`
	Status      ServerStatus      `json:"status"`
	Tools       []Tool            `json:"tools"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// ServerType 表示服务器类型
type ServerType string

const (
	ServerTypeLocal  ServerType = "local"
	ServerTypeRemote ServerType = "remote"
)

// ServerStatus 表示服务器状态
type ServerStatus string

const (
	ServerStatusRunning ServerStatus = "running"
	ServerStatusStopped ServerStatus = "stopped"
	ServerStatusError   ServerStatus = "error"
)

// ToolResult 表示工具执行结果
type ToolResult struct {
	ToolID    string      `json:"tool_id"`
	Status    string      `json:"status"`
	Result    interface{} `json:"result,omitempty"`
	Error     string      `json:"error,omitempty"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
}

// ToolExecutionStatus 表示工具执行状态
type ToolExecutionStatus string

const (
	ToolExecutionStatusSuccess ToolExecutionStatus = "success"
	ToolExecutionStatusError   ToolExecutionStatus = "error"
	ToolExecutionStatusTimeout ToolExecutionStatus = "timeout"
)

// ToolExecutionResult 表示工具执行结果
type ToolExecutionResult struct {
	Status    ToolExecutionStatus `json:"status"`
	Result    interface{}         `json:"result,omitempty"`
	Error     string              `json:"error,omitempty"`
	StartTime time.Time           `json:"start_time"`
	EndTime   time.Time           `json:"end_time"`
}

// ToolExecutor 定义了工具执行器接口
type ToolExecutor interface {
	Execute(ctx context.Context, tool *Tool, params map[string]interface{}) (*ToolExecutionResult, error)
}

// ToolManager 定义了工具管理接口
type ToolManager interface {
	// 服务器管理
	AddServer(ctx context.Context, server *Server) error
	RemoveServer(ctx context.Context, serverID string) error
	GetServer(ctx context.Context, serverID string) (*Server, error)
	ListServers(ctx context.Context) ([]*Server, error)
	StartServer(ctx context.Context, serverID string) error
	StopServer(ctx context.Context, serverID string) error

	// 工具管理
	GetTool(ctx context.Context, toolID string) (*Tool, error)
	ListTools(ctx context.Context) ([]*Tool, error)
	ExecuteTool(ctx context.Context, toolID string, params map[string]interface{}) (*ToolResult, error)

	// 市场相关
	SearchTools(ctx context.Context, query string) ([]*Tool, error)
	DownloadTool(ctx context.Context, toolID string) error
	UpdateTool(ctx context.Context, toolID string) error

	// 配置管理
	SetAutoApprove(ctx context.Context, enabled bool) error
	GetAutoApprove(ctx context.Context) bool
	SetTimeout(ctx context.Context, timeout time.Duration) error
	GetTimeout(ctx context.Context) time.Duration
}

// ValidateParameters 验证工具参数
func (t *Tool) ValidateParameters(params map[string]interface{}) error {
	// 检查必需参数
	for _, param := range t.Parameters {
		if param.Required {
			if _, exists := params[param.Name]; !exists {
				return fmt.Errorf("missing required parameter: %s", param.Name)
			}
		}
	}

	// 检查参数类型
	for name, value := range params {
		// 查找参数定义
		var paramDef *ToolParameter
		for _, p := range t.Parameters {
			if p.Name == name {
				paramDef = &p
				break
			}
		}

		if paramDef == nil {
			return fmt.Errorf("unknown parameter: %s", name)
		}

		// 验证参数类型
		if err := validateParameterType(paramDef.Type, value); err != nil {
			return fmt.Errorf("invalid type for parameter %s: %v", name, err)
		}
	}

	return nil
}

// validateParameterType 验证参数类型
func validateParameterType(paramType string, value interface{}) error {
	switch paramType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "number":
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			// 数字类型都接受
		default:
			return fmt.Errorf("expected number, got %T", value)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	default:
		return fmt.Errorf("unsupported parameter type: %s", paramType)
	}
	return nil
}
