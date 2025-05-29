package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/liangsj/vimcoplit/internal/models"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// 测试服务器配置
	if cfg.Server.Host != "localhost" {
		t.Errorf("expected server host to be 'localhost', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server port to be 8080, got %d", cfg.Server.Port)
	}

	// 测试模型配置
	if cfg.Model.Type != models.ModelTypeClaude {
		t.Errorf("expected model type to be %s, got %s", models.ModelTypeClaude, cfg.Model.Type)
	}
	if cfg.Model.MaxTokens != 4096 {
		t.Errorf("expected max tokens to be 4096, got %d", cfg.Model.MaxTokens)
	}
	if cfg.Model.Temperature != 0.7 {
		t.Errorf("expected temperature to be 0.7, got %f", cfg.Model.Temperature)
	}

	// 测试日志配置
	if cfg.Log.Level != "info" {
		t.Errorf("expected log level to be 'info', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.File != "vimcoplit.log" {
		t.Errorf("expected log file to be 'vimcoplit.log', got '%s'", cfg.Log.File)
	}

	// 测试文件配置
	if cfg.File.MaxFileSize != 10*1024*1024 {
		t.Errorf("expected max file size to be 10MB, got %d", cfg.File.MaxFileSize)
	}
	expectedExts := []string{".go", ".lua", ".md", ".txt"}
	if len(cfg.File.AllowedExts) != len(expectedExts) {
		t.Errorf("expected %d allowed extensions, got %d", len(expectedExts), len(cfg.File.AllowedExts))
	}

	// 测试命令配置
	if cfg.Command.Timeout != 30 {
		t.Errorf("expected command timeout to be 30, got %d", cfg.Command.Timeout)
	}
	expectedCmds := []string{"git", "go", "nvim"}
	if len(cfg.Command.AllowedCmds) != len(expectedCmds) {
		t.Errorf("expected %d allowed commands, got %d", len(expectedCmds), len(cfg.Command.AllowedCmds))
	}
}

func TestLoadConfig(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "vimcoplit-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试默认配置加载
	cfg, err := LoadConfig("")
	if err != nil {
		t.Errorf("failed to load default config: %v", err)
	}
	if cfg == nil {
		t.Error("expected config to be non-nil")
	}

	// 测试从文件加载配置
	configPath := filepath.Join(tempDir, "config.json")
	cfg, err = LoadConfig(configPath)
	if err != nil {
		t.Errorf("failed to load config from file: %v", err)
	}
	if cfg == nil {
		t.Error("expected config to be non-nil")
	}

	// 验证配置文件是否被创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestSaveConfig(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "vimcoplit-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试配置
	cfg := DefaultConfig()
	cfg.Server.Port = 9090
	cfg.Model.Type = models.ModelTypeDoubao

	// 保存配置
	configPath := filepath.Join(tempDir, "config.json")
	err = SaveConfig(configPath, cfg)
	if err != nil {
		t.Errorf("failed to save config: %v", err)
	}

	// 验证文件是否被创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}

	// 重新加载配置并验证
	loadedCfg, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("failed to load saved config: %v", err)
	}
	if loadedCfg.Server.Port != 9090 {
		t.Errorf("expected port to be 9090, got %d", loadedCfg.Server.Port)
	}
	if loadedCfg.Model.Type != models.ModelTypeDoubao {
		t.Errorf("expected model type to be %s, got %s", models.ModelTypeDoubao, loadedCfg.Model.Type)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// 设置环境变量
	os.Setenv("VIMCOPLIT_HOST", "test-host")
	os.Setenv("VIMCOPLIT_PORT", "9090")
	os.Setenv("VIMCOPLIT_MODEL_TYPE", "doubao")
	os.Setenv("VIMCOPLIT_API_KEY", "test-key")
	os.Setenv("VIMCOPLIT_MAX_TOKENS", "2048")
	os.Setenv("VIMCOPLIT_TEMPERATURE", "0.5")
	os.Setenv("VIMCOPLIT_LOG_LEVEL", "debug")
	os.Setenv("VIMCOPLIT_LOG_FILE", "test.log")

	// 创建配置并加载环境变量
	cfg := DefaultConfig()
	loadFromEnv(cfg)

	// 验证环境变量覆盖
	if cfg.Server.Host != "test-host" {
		t.Errorf("expected host to be 'test-host', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected port to be 9090, got %d", cfg.Server.Port)
	}
	if cfg.Model.Type != models.ModelTypeDoubao {
		t.Errorf("expected model type to be %s, got %s", models.ModelTypeDoubao, cfg.Model.Type)
	}
	if cfg.Model.APIKey != "test-key" {
		t.Errorf("expected API key to be 'test-key', got '%s'", cfg.Model.APIKey)
	}
	if cfg.Model.MaxTokens != 2048 {
		t.Errorf("expected max tokens to be 2048, got %d", cfg.Model.MaxTokens)
	}
	if cfg.Model.Temperature != 0.5 {
		t.Errorf("expected temperature to be 0.5, got %f", cfg.Model.Temperature)
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("expected log level to be 'debug', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.File != "test.log" {
		t.Errorf("expected log file to be 'test.log', got '%s'", cfg.Log.File)
	}

	// 清理环境变量
	os.Unsetenv("VIMCOPLIT_HOST")
	os.Unsetenv("VIMCOPLIT_PORT")
	os.Unsetenv("VIMCOPLIT_MODEL_TYPE")
	os.Unsetenv("VIMCOPLIT_API_KEY")
	os.Unsetenv("VIMCOPLIT_MAX_TOKENS")
	os.Unsetenv("VIMCOPLIT_TEMPERATURE")
	os.Unsetenv("VIMCOPLIT_LOG_LEVEL")
	os.Unsetenv("VIMCOPLIT_LOG_FILE")
}

func TestGetConfig(t *testing.T) {
	// 测试获取默认配置
	cfg := GetConfig()
	if cfg == nil {
		t.Error("expected config to be non-nil")
	}

	// 修改配置
	cfg.Server.Port = 9090

	// 再次获取配置，应该返回相同的实例
	cfg2 := GetConfig()
	if cfg2.Server.Port != 9090 {
		t.Errorf("expected port to be 9090, got %d", cfg2.Server.Port)
	}
}
