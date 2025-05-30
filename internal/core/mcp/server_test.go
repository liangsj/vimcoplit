package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalServerRunner(t *testing.T) {
	// 创建一个测试服务器
	server := &Server{
		ID:   "test-local-server",
		Name: "Test Local Server",
		Type: ServerTypeLocal,
		Metadata: map[string]string{
			"start_cmd": "echo 'Server started'",
			"work_dir":  ".",
			"env":       "TEST_ENV=test",
		},
	}

	// 创建运行器
	runner := NewLocalServerRunner(server)

	// 测试启动
	ctx := context.Background()
	if err := runner.Start(ctx); err != nil {
		t.Errorf("Failed to start server: %v", err)
	}

	// 验证状态
	if status := runner.Status(); status != ServerStatusRunning {
		t.Errorf("Expected status %s, got %s", ServerStatusRunning, status)
	}

	// 测试停止
	if err := runner.Stop(ctx); err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}

	// 验证状态
	if status := runner.Status(); status != ServerStatusStopped {
		t.Errorf("Expected status %s, got %s", ServerStatusStopped, status)
	}
}

func TestRemoteServerRunner(t *testing.T) {
	// 创建一个测试 HTTP 服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	// 创建一个测试服务器
	server := &Server{
		ID:   "test-remote-server",
		Name: "Test Remote Server",
		Type: ServerTypeRemote,
		URL:  ts.URL,
	}

	// 创建运行器
	runner := NewRemoteServerRunner(server)

	// 测试启动
	ctx := context.Background()
	if err := runner.Start(ctx); err != nil {
		t.Errorf("Failed to start server: %v", err)
	}

	// 验证状态
	if status := runner.Status(); status != ServerStatusRunning {
		t.Errorf("Expected status %s, got %s", ServerStatusRunning, status)
	}

	// 测试健康检查
	if err := runner.HealthCheck(ctx); err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// 测试停止
	if err := runner.Stop(ctx); err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}

	// 验证状态
	if status := runner.Status(); status != ServerStatusStopped {
		t.Errorf("Expected status %s, got %s", ServerStatusStopped, status)
	}
}

func TestManagerServerOperations(t *testing.T) {
	// 创建一个测试管理器
	manager := NewManager("test_config.json")

	// 创建一个测试服务器
	server := &Server{
		ID:   "test-server",
		Name: "Test Server",
		Type: ServerTypeLocal,
		Metadata: map[string]string{
			"start_cmd": "echo 'Server started'",
		},
	}

	// 添加服务器
	ctx := context.Background()
	if err := manager.AddServer(ctx, server); err != nil {
		t.Errorf("Failed to add server: %v", err)
	}

	// 测试启动服务器
	if err := manager.StartServer(ctx, server.ID); err != nil {
		t.Errorf("Failed to start server: %v", err)
	}

	// 验证服务器状态
	updatedServer, err := manager.GetServer(ctx, server.ID)
	if err != nil {
		t.Errorf("Failed to get server: %v", err)
	}
	if updatedServer.Status != ServerStatusRunning {
		t.Errorf("Expected status %s, got %s", ServerStatusRunning, updatedServer.Status)
	}

	// 测试停止服务器
	if err := manager.StopServer(ctx, server.ID); err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}

	// 验证服务器状态
	updatedServer, err = manager.GetServer(ctx, server.ID)
	if err != nil {
		t.Errorf("Failed to get server: %v", err)
	}
	if updatedServer.Status != ServerStatusStopped {
		t.Errorf("Expected status %s, got %s", ServerStatusStopped, updatedServer.Status)
	}
}
