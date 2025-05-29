package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager 是 ToolManager 接口的具体实现
type Manager struct {
	servers     map[string]*Server
	tools       map[string]*Tool
	autoApprove bool
	timeout     time.Duration
	mu          sync.RWMutex
	configPath  string
	executors   map[string]ToolExecutor
}

// NewManager 创建一个新的工具管理器
func NewManager(configPath string) *Manager {
	return &Manager{
		servers:     make(map[string]*Server),
		tools:       make(map[string]*Tool),
		autoApprove: false,
		timeout:     30 * time.Second,
		configPath:  configPath,
		executors:   make(map[string]ToolExecutor),
	}
}

// AddServer 添加一个新的 MCP 服务器
func (m *Manager) AddServer(ctx context.Context, server *Server) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if server.ID == "" {
		server.ID = uuid.New().String()
	}
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	m.servers[server.ID] = server
	return m.saveConfig()
}

// RemoveServer 移除一个 MCP 服务器
func (m *Manager) RemoveServer(ctx context.Context, serverID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[serverID]; !exists {
		return errors.New("server not found")
	}

	// 移除服务器相关的所有工具
	for toolID, tool := range m.tools {
		if tool.ServerID == serverID {
			delete(m.tools, toolID)
		}
	}

	delete(m.servers, serverID)
	return m.saveConfig()
}

// GetServer 获取服务器信息
func (m *Manager) GetServer(ctx context.Context, serverID string) (*Server, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	server, exists := m.servers[serverID]
	if !exists {
		return nil, errors.New("server not found")
	}
	return server, nil
}

// ListServers 列出所有服务器
func (m *Manager) ListServers(ctx context.Context) ([]*Server, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	servers := make([]*Server, 0, len(m.servers))
	for _, server := range m.servers {
		servers = append(servers, server)
	}
	return servers, nil
}

// StartServer 启动服务器
func (m *Manager) StartServer(ctx context.Context, serverID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	server, exists := m.servers[serverID]
	if !exists {
		return errors.New("server not found")
	}

	// TODO: 实现实际的服务器启动逻辑
	server.Status = ServerStatusRunning
	server.UpdatedAt = time.Now()
	return m.saveConfig()
}

// StopServer 停止服务器
func (m *Manager) StopServer(ctx context.Context, serverID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	server, exists := m.servers[serverID]
	if !exists {
		return errors.New("server not found")
	}

	// TODO: 实现实际的服务器停止逻辑
	server.Status = ServerStatusStopped
	server.UpdatedAt = time.Now()
	return m.saveConfig()
}

// GetTool 获取工具信息
func (m *Manager) GetTool(ctx context.Context, toolID string) (*Tool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tool, exists := m.tools[toolID]
	if !exists {
		return nil, errors.New("tool not found")
	}
	return tool, nil
}

// ListTools 列出所有工具
func (m *Manager) ListTools(ctx context.Context) ([]*Tool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tools := make([]*Tool, 0, len(m.tools))
	for _, tool := range m.tools {
		tools = append(tools, tool)
	}
	return tools, nil
}

// ExecuteTool 执行工具
func (m *Manager) ExecuteTool(ctx context.Context, toolID string, params map[string]interface{}) (*ToolResult, error) {
	m.mu.RLock()
	tool, exists := m.tools[toolID]
	m.mu.RUnlock()

	if !exists {
		return nil, errors.New("tool not found")
	}

	// 检查服务器状态
	server, err := m.GetServer(ctx, tool.ServerID)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}

	if server.Status != ServerStatusRunning {
		return nil, errors.New("server is not running")
	}

	// 获取执行器
	executor, exists := m.executors[tool.ServerID]
	if !exists {
		// 根据服务器类型创建执行器
		switch server.Type {
		case ServerTypeLocal:
			executor = NewLocalExecutor()
		case ServerTypeRemote:
			executor = NewHTTPExecutor(m.timeout)
		default:
			return nil, fmt.Errorf("unsupported server type: %s", server.Type)
		}
		m.executors[tool.ServerID] = executor
	}

	// 执行工具
	result, err := executor.Execute(ctx, tool, params)
	if err != nil {
		return nil, err
	}

	// 转换结果
	return &ToolResult{
		ToolID:    toolID,
		Status:    string(result.Status),
		Result:    result.Result,
		Error:     result.Error,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
	}, nil
}

// SearchTools 搜索工具
func (m *Manager) SearchTools(ctx context.Context, query string) ([]*Tool, error) {
	// TODO: 实现工具市场搜索
	return nil, nil
}

// DownloadTool 下载工具
func (m *Manager) DownloadTool(ctx context.Context, toolID string) error {
	// TODO: 实现工具下载
	return nil
}

// UpdateTool 更新工具
func (m *Manager) UpdateTool(ctx context.Context, toolID string) error {
	// TODO: 实现工具更新
	return nil
}

// SetAutoApprove 设置自动审批
func (m *Manager) SetAutoApprove(ctx context.Context, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.autoApprove = enabled
	return m.saveConfig()
}

// GetAutoApprove 获取自动审批状态
func (m *Manager) GetAutoApprove(ctx context.Context) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.autoApprove
}

// SetTimeout 设置超时时间
func (m *Manager) SetTimeout(ctx context.Context, timeout time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.timeout = timeout
	return m.saveConfig()
}

// GetTimeout 获取超时时间
func (m *Manager) GetTimeout(ctx context.Context) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.timeout
}

// RegisterLocalTool 注册本地工具
func (m *Manager) RegisterLocalTool(serverID string, tool *Tool, handler ToolHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查服务器
	server, exists := m.servers[serverID]
	if !exists {
		return errors.New("server not found")
	}

	if server.Type != ServerTypeLocal {
		return errors.New("server is not a local server")
	}

	// 注册工具
	tool.ServerID = serverID
	tool.CreatedAt = time.Now()
	tool.UpdatedAt = time.Now()
	m.tools[tool.ID] = tool

	// 注册处理函数
	executor, exists := m.executors[serverID]
	if !exists {
		executor = NewLocalExecutor()
		m.executors[serverID] = executor
	}

	localExecutor, ok := executor.(*LocalExecutor)
	if !ok {
		return errors.New("invalid executor type")
	}

	localExecutor.RegisterHandler(tool.ID, handler)
	return m.saveConfig()
}

// saveConfig 保存配置到文件
func (m *Manager) saveConfig() error {
	config := struct {
		Servers     map[string]*Server `json:"servers"`
		Tools       map[string]*Tool   `json:"tools"`
		AutoApprove bool               `json:"auto_approve"`
		Timeout     time.Duration      `json:"timeout"`
	}{
		Servers:     m.servers,
		Tools:       m.tools,
		AutoApprove: m.autoApprove,
		Timeout:     m.timeout,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(m.configPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// loadConfig 从文件加载配置
func (m *Manager) loadConfig() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var config struct {
		Servers     map[string]*Server `json:"servers"`
		Tools       map[string]*Tool   `json:"tools"`
		AutoApprove bool               `json:"auto_approve"`
		Timeout     time.Duration      `json:"timeout"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.servers = config.Servers
	m.tools = config.Tools
	m.autoApprove = config.AutoApprove
	m.timeout = config.Timeout

	return nil
}

// LoadToolFromConfig 从配置文件加载工具
func (m *Manager) LoadToolFromConfig(ctx context.Context, configPath string) error {
	config, err := LoadToolConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load tool config: %v", err)
	}

	// 创建工具
	tool := &Tool{
		ID:          config.ID,
		Name:        config.Name,
		Description: config.Description,
		Version:     config.Version,
		Author:      config.Author,
		Parameters:  config.Parameters,
		Metadata:    config.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 添加工具
	m.mu.Lock()
	m.tools[tool.ID] = tool
	m.mu.Unlock()

	return m.saveConfig()
}

// LoadServerFromConfig 从配置文件加载服务器
func (m *Manager) LoadServerFromConfig(ctx context.Context, configPath string) error {
	config, err := LoadServerConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load server config: %v", err)
	}

	// 创建服务器
	server := &Server{
		ID:          config.ID,
		Name:        config.Name,
		Description: config.Description,
		Version:     config.Version,
		Type:        config.Type,
		URL:         config.URL,
		Status:      ServerStatusStopped,
		Metadata:    config.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 添加服务器
	m.mu.Lock()
	m.servers[server.ID] = server
	m.mu.Unlock()

	// 加载工具
	for _, toolConfig := range config.Tools {
		tool := &Tool{
			ID:          toolConfig.ID,
			Name:        toolConfig.Name,
			Description: toolConfig.Description,
			Version:     toolConfig.Version,
			Author:      toolConfig.Author,
			Parameters:  toolConfig.Parameters,
			ServerID:    server.ID,
			Metadata:    toolConfig.Metadata,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		m.mu.Lock()
		m.tools[tool.ID] = tool
		m.mu.Unlock()
	}

	return m.saveConfig()
}

// LoadConfigsFromDirectory 从目录加载所有配置文件
func (m *Manager) LoadConfigsFromDirectory(ctx context.Context, dirPath string) error {
	// 加载服务器配置
	serverConfigs, err := filepath.Glob(filepath.Join(dirPath, "servers", "*.json"))
	if err != nil {
		return fmt.Errorf("failed to find server configs: %v", err)
	}

	for _, configPath := range serverConfigs {
		if err := m.LoadServerFromConfig(ctx, configPath); err != nil {
			return fmt.Errorf("failed to load server config %s: %v", configPath, err)
		}
	}

	// 加载工具配置
	toolConfigs, err := filepath.Glob(filepath.Join(dirPath, "tools", "*.json"))
	if err != nil {
		return fmt.Errorf("failed to find tool configs: %v", err)
	}

	for _, configPath := range toolConfigs {
		if err := m.LoadToolFromConfig(ctx, configPath); err != nil {
			return fmt.Errorf("failed to load tool config %s: %v", configPath, err)
		}
	}

	return nil
}
