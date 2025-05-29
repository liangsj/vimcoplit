package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/liangsj/vimcoplit/internal/models"
)

// Config 定义了应用程序的配置结构
type Config struct {
	// 服务器配置
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`

	// AI模型配置
	Model struct {
		Type        models.ModelType `json:"type"`
		APIKey      string           `json:"api_key"`
		MaxTokens   int              `json:"max_tokens"`
		Temperature float64          `json:"temperature"`
	} `json:"model"`

	// 日志配置
	Log struct {
		Level      string `json:"level"`
		File       string `json:"file"`
		MaxSize    int    `json:"max_size"`
		MaxBackups int    `json:"max_backups"`
		MaxAge     int    `json:"max_age"`
	} `json:"log"`

	// 文件操作配置
	File struct {
		MaxFileSize int64    `json:"max_file_size"`
		AllowedExts []string `json:"allowed_exts"`
	} `json:"file"`

	// 命令执行配置
	Command struct {
		Timeout     int      `json:"timeout"`
		AllowedCmds []string `json:"allowed_cmds"`
	} `json:"command"`
}

var (
	config *Config
	once   sync.Once
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}{
			Host: "localhost",
			Port: 8080,
		},
		Model: struct {
			Type        models.ModelType `json:"type"`
			APIKey      string           `json:"api_key"`
			MaxTokens   int              `json:"max_tokens"`
			Temperature float64          `json:"temperature"`
		}{
			Type:        models.ModelTypeClaude,
			MaxTokens:   4096,
			Temperature: 0.7,
		},
		Log: struct {
			Level      string `json:"level"`
			File       string `json:"file"`
			MaxSize    int    `json:"max_size"`
			MaxBackups int    `json:"max_backups"`
			MaxAge     int    `json:"max_age"`
		}{
			Level:      "info",
			File:       "vimcoplit.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
		},
		File: struct {
			MaxFileSize int64    `json:"max_file_size"`
			AllowedExts []string `json:"allowed_exts"`
		}{
			MaxFileSize: 10 * 1024 * 1024, // 10MB
			AllowedExts: []string{".go", ".lua", ".md", ".txt"},
		},
		Command: struct {
			Timeout     int      `json:"timeout"`
			AllowedCmds []string `json:"allowed_cmds"`
		}{
			Timeout:     30,
			AllowedCmds: []string{"git", "go", "nvim"},
		},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	once.Do(func() {
		config = DefaultConfig()
	})

	// 如果配置文件路径为空，使用默认路径
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			configPath = "config.json"
		} else {
			configPath = filepath.Join(homeDir, ".vimcoplit", "config.json")
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，创建默认配置文件
			if err := SaveConfig(configPath, config); err != nil {
				return nil, fmt.Errorf("failed to create default config: %v", err)
			}
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析配置文件
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// 从环境变量加载配置
	loadFromEnv(config)

	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(configPath string, cfg *Config) error {
	// 确保配置目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 序列化配置
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetConfig 返回当前配置
func GetConfig() *Config {
	if config == nil {
		config = DefaultConfig()
	}
	return config
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(cfg *Config) {
	// 服务器配置
	if host := os.Getenv("VIMCOPLIT_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("VIMCOPLIT_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Server.Port)
	}

	// 模型配置
	if modelType := os.Getenv("VIMCOPLIT_MODEL_TYPE"); modelType != "" {
		cfg.Model.Type = models.ModelType(modelType)
	}
	if apiKey := os.Getenv("VIMCOPLIT_API_KEY"); apiKey != "" {
		cfg.Model.APIKey = apiKey
	}
	if maxTokens := os.Getenv("VIMCOPLIT_MAX_TOKENS"); maxTokens != "" {
		fmt.Sscanf(maxTokens, "%d", &cfg.Model.MaxTokens)
	}
	if temp := os.Getenv("VIMCOPLIT_TEMPERATURE"); temp != "" {
		fmt.Sscanf(temp, "%f", &cfg.Model.Temperature)
	}

	// 日志配置
	if level := os.Getenv("VIMCOPLIT_LOG_LEVEL"); level != "" {
		cfg.Log.Level = level
	}
	if file := os.Getenv("VIMCOPLIT_LOG_FILE"); file != "" {
		cfg.Log.File = file
	}
}
