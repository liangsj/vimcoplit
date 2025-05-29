package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ToolConfig 表示工具配置
type ToolConfig struct {
	// 工具基本信息
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Parameters  []ToolParameter   `json:"parameters"`
	Metadata    map[string]string `json:"metadata"`

	// 执行配置
	Timeout     int64 `json:"timeout,omitempty"`     // 超时时间（秒）
	RetryCount  int   `json:"retry_count,omitempty"` // 重试次数
	RetryDelay  int64 `json:"retry_delay,omitempty"` // 重试延迟（秒）
	Concurrency int   `json:"concurrency,omitempty"` // 并发数
	RateLimit   int   `json:"rate_limit,omitempty"`  // 速率限制（每秒请求数）

	// 安全配置
	RequireAuth bool     `json:"require_auth,omitempty"` // 是否需要认证
	AllowRoles  []string `json:"allow_roles,omitempty"`  // 允许的角色
	AllowIPs    []string `json:"allow_ips,omitempty"`    // 允许的 IP 地址

	// 日志配置
	LogLevel    string `json:"log_level,omitempty"`     // 日志级别
	LogFile     string `json:"log_file,omitempty"`      // 日志文件
	LogFormat   string `json:"log_format,omitempty"`    // 日志格式
	LogMaxSize  int    `json:"log_max_size,omitempty"`  // 日志文件最大大小（MB）
	LogMaxFiles int    `json:"log_max_files,omitempty"` // 最大日志文件数
}

// ServerConfig 表示服务器配置
type ServerConfig struct {
	// 服务器基本信息
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Version     string            `json:"version"`
	Type        ServerType        `json:"type"`
	URL         string            `json:"url,omitempty"`
	Metadata    map[string]string `json:"metadata"`

	// 工具配置
	Tools []ToolConfig `json:"tools"`

	// 服务器配置
	Port            int      `json:"port,omitempty"`             // 服务器端口
	Host            string   `json:"host,omitempty"`             // 服务器主机
	SSLEnabled      bool     `json:"ssl_enabled,omitempty"`      // 是否启用 SSL
	SSLCertFile     string   `json:"ssl_cert_file,omitempty"`    // SSL 证书文件
	SSLKeyFile      string   `json:"ssl_key_file,omitempty"`     // SSL 密钥文件
	AllowedOrigins  []string `json:"allowed_origins,omitempty"`  // 允许的源
	AllowedMethods  []string `json:"allowed_methods,omitempty"`  // 允许的方法
	AllowedHeaders  []string `json:"allowed_headers,omitempty"`  // 允许的头部
	MaxRequestSize  int64    `json:"max_request_size,omitempty"` // 最大请求大小（字节）
	ReadTimeout     int64    `json:"read_timeout,omitempty"`     // 读取超时（秒）
	WriteTimeout    int64    `json:"write_timeout,omitempty"`    // 写入超时（秒）
	IdleTimeout     int64    `json:"idle_timeout,omitempty"`     // 空闲超时（秒）
	ShutdownTimeout int64    `json:"shutdown_timeout,omitempty"` // 关闭超时（秒）
}

// LoadToolConfig 从文件加载工具配置
func LoadToolConfig(path string) (*ToolConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read tool config: %v", err)
	}

	var config ToolConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse tool config: %v", err)
	}

	if err := validateToolConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid tool config: %v", err)
	}

	return &config, nil
}

// LoadServerConfig 从文件加载服务器配置
func LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read server config: %v", err)
	}

	var config ServerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse server config: %v", err)
	}

	if err := validateServerConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid server config: %v", err)
	}

	return &config, nil
}

// SaveToolConfig 保存工具配置到文件
func SaveToolConfig(config *ToolConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tool config: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write tool config: %v", err)
	}

	return nil
}

// SaveServerConfig 保存服务器配置到文件
func SaveServerConfig(config *ServerConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal server config: %v", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write server config: %v", err)
	}

	return nil
}

// validateToolConfig 验证工具配置
func validateToolConfig(config *ToolConfig) error {
	if config.ID == "" {
		return fmt.Errorf("tool ID is required")
	}
	if config.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	if config.Version == "" {
		return fmt.Errorf("tool version is required")
	}

	// 验证参数配置
	for i, param := range config.Parameters {
		if param.Name == "" {
			return fmt.Errorf("parameter name is required at index %d", i)
		}
		if param.Type == "" {
			return fmt.Errorf("parameter type is required for parameter %s", param.Name)
		}
	}

	// 验证超时配置
	if config.Timeout < 0 {
		return fmt.Errorf("timeout must be non-negative")
	}

	// 验证重试配置
	if config.RetryCount < 0 {
		return fmt.Errorf("retry count must be non-negative")
	}
	if config.RetryDelay < 0 {
		return fmt.Errorf("retry delay must be non-negative")
	}

	// 验证并发配置
	if config.Concurrency < 0 {
		return fmt.Errorf("concurrency must be non-negative")
	}

	// 验证速率限制
	if config.RateLimit < 0 {
		return fmt.Errorf("rate limit must be non-negative")
	}

	// 验证日志配置
	if config.LogMaxSize < 0 {
		return fmt.Errorf("log max size must be non-negative")
	}
	if config.LogMaxFiles < 0 {
		return fmt.Errorf("log max files must be non-negative")
	}

	return nil
}

// validateServerConfig 验证服务器配置
func validateServerConfig(config *ServerConfig) error {
	if config.ID == "" {
		return fmt.Errorf("server ID is required")
	}
	if config.Name == "" {
		return fmt.Errorf("server name is required")
	}
	if config.Version == "" {
		return fmt.Errorf("server version is required")
	}
	if config.Type == "" {
		return fmt.Errorf("server type is required")
	}

	// 验证工具配置
	for i, tool := range config.Tools {
		if err := validateToolConfig(&tool); err != nil {
			return fmt.Errorf("invalid tool config at index %d: %v", i, err)
		}
	}

	// 验证端口配置
	if config.Port < 0 || config.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}

	// 验证 SSL 配置
	if config.SSLEnabled {
		if config.SSLCertFile == "" {
			return fmt.Errorf("SSL certificate file is required when SSL is enabled")
		}
		if config.SSLKeyFile == "" {
			return fmt.Errorf("SSL key file is required when SSL is enabled")
		}
	}

	// 验证超时配置
	if config.ReadTimeout < 0 {
		return fmt.Errorf("read timeout must be non-negative")
	}
	if config.WriteTimeout < 0 {
		return fmt.Errorf("write timeout must be non-negative")
	}
	if config.IdleTimeout < 0 {
		return fmt.Errorf("idle timeout must be non-negative")
	}
	if config.ShutdownTimeout < 0 {
		return fmt.Errorf("shutdown timeout must be non-negative")
	}

	return nil
}
