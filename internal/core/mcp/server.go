package mcp

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

// ServerRunner 定义了服务器运行器接口
type ServerRunner interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status() ServerStatus
	HealthCheck(ctx context.Context) error
}

// LocalServerRunner 是本地服务器的运行器
type LocalServerRunner struct {
	server     *Server
	cmd        *exec.Cmd
	mu         sync.RWMutex
	status     ServerStatus
	stopChan   chan struct{}
	healthURL  string
	httpClient *http.Client
}

// NewLocalServerRunner 创建一个新的本地服务器运行器
func NewLocalServerRunner(server *Server) *LocalServerRunner {
	return &LocalServerRunner{
		server:     server,
		status:     ServerStatusStopped,
		stopChan:   make(chan struct{}),
		healthURL:  server.Metadata["health_url"],
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Start 启动本地服务器
func (r *LocalServerRunner) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status == ServerStatusRunning {
		return nil
	}

	// 获取启动命令
	cmd := r.server.Metadata["start_cmd"]
	if cmd == "" {
		return fmt.Errorf("no start command specified for server %s", r.server.ID)
	}

	// 创建命令
	r.cmd = exec.CommandContext(ctx, "sh", "-c", cmd)

	// 设置工作目录
	if workDir := r.server.Metadata["work_dir"]; workDir != "" {
		r.cmd.Dir = workDir
	}

	// 设置环境变量
	if env := r.server.Metadata["env"]; env != "" {
		r.cmd.Env = append(r.cmd.Env, env)
	}

	// 启动进程
	if err := r.cmd.Start(); err != nil {
		r.status = ServerStatusError
		return fmt.Errorf("failed to start server: %v", err)
	}

	// 更新状态
	r.status = ServerStatusRunning

	// 启动健康检查
	go r.healthCheck()

	return nil
}

// Stop 停止本地服务器
func (r *LocalServerRunner) Stop(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != ServerStatusRunning {
		return nil
	}

	// 发送停止信号
	close(r.stopChan)

	// 停止进程
	if r.cmd != nil && r.cmd.Process != nil {
		if err := r.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop server: %v", err)
		}
	}

	// 更新状态
	r.status = ServerStatusStopped
	return nil
}

// Status 获取服务器状态
func (r *LocalServerRunner) Status() ServerStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.status
}

// HealthCheck 执行健康检查
func (r *LocalServerRunner) HealthCheck(ctx context.Context) error {
	if r.healthURL == "" {
		return nil
	}

	resp, err := r.httpClient.Get(r.healthURL)
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// healthCheck 定期执行健康检查
func (r *LocalServerRunner) healthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.HealthCheck(context.Background()); err != nil {
				r.mu.Lock()
				r.status = ServerStatusError
				r.mu.Unlock()
			}
		case <-r.stopChan:
			return
		}
	}
}

// RemoteServerRunner 是远程服务器的运行器
type RemoteServerRunner struct {
	server     *Server
	mu         sync.RWMutex
	status     ServerStatus
	httpClient *http.Client
}

// NewRemoteServerRunner 创建一个新的远程服务器运行器
func NewRemoteServerRunner(server *Server) *RemoteServerRunner {
	return &RemoteServerRunner{
		server:     server,
		status:     ServerStatusStopped,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Start 启动远程服务器
func (r *RemoteServerRunner) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status == ServerStatusRunning {
		return nil
	}

	// 检查服务器是否可访问
	if err := r.HealthCheck(ctx); err != nil {
		r.status = ServerStatusError
		return fmt.Errorf("server is not accessible: %v", err)
	}

	// 更新状态
	r.status = ServerStatusRunning
	return nil
}

// Stop 停止远程服务器
func (r *RemoteServerRunner) Stop(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != ServerStatusRunning {
		return nil
	}

	// 更新状态
	r.status = ServerStatusStopped
	return nil
}

// Status 获取服务器状态
func (r *RemoteServerRunner) Status() ServerStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.status
}

// HealthCheck 执行健康检查
func (r *RemoteServerRunner) HealthCheck(ctx context.Context) error {
	healthURL := r.server.URL + "/health"
	if customHealthURL := r.server.Metadata["health_url"]; customHealthURL != "" {
		healthURL = customHealthURL
	}

	resp, err := r.httpClient.Get(healthURL)
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
