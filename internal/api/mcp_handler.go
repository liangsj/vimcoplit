package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/liangsj/vimcoplit/internal/core/mcp"
)

// MCPHandler 处理 MCP 相关的 HTTP 请求
type MCPHandler struct {
	manager *mcp.Manager
}

// NewMCPHandler 创建一个新的 MCP 处理器
func NewMCPHandler(manager *mcp.Manager) *MCPHandler {
	return &MCPHandler{
		manager: manager,
	}
}

// RegisterRoutes 注册路由
func (h *MCPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/mcp/servers", h.handleServers)
	mux.HandleFunc("/api/mcp/tools", h.handleTools)
	mux.HandleFunc("/api/mcp/config", h.handleConfig)
}

// handleServers 处理服务器相关的请求
func (h *MCPHandler) handleServers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listServers(w, r)
	case http.MethodPost:
		h.addServer(w, r)
	case http.MethodDelete:
		h.removeServer(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTools 处理工具相关的请求
func (h *MCPHandler) handleTools(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTools(w, r)
	case http.MethodPost:
		h.executeTool(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleConfig 处理配置相关的请求
func (h *MCPHandler) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getConfig(w, r)
	case http.MethodPut:
		h.updateConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listServers 列出所有服务器
func (h *MCPHandler) listServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.manager.ListServers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(servers)
}

// addServer 添加新服务器
func (h *MCPHandler) addServer(w http.ResponseWriter, r *http.Request) {
	var server mcp.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.manager.AddServer(r.Context(), &server); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(server)
}

// removeServer 移除服务器
func (h *MCPHandler) removeServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("id")
	if serverID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	if err := h.manager.RemoveServer(r.Context(), serverID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listTools 列出所有工具
func (h *MCPHandler) listTools(w http.ResponseWriter, r *http.Request) {
	tools, err := h.manager.ListTools(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tools)
}

// executeTool 执行工具
func (h *MCPHandler) executeTool(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ToolID  string                 `json:"tool_id"`
		Params  map[string]interface{} `json:"params"`
		Timeout time.Duration          `json:"timeout,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 如果请求指定了超时，使用请求的超时
	ctx := r.Context()
	if req.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, req.Timeout)
		defer cancel()
	}

	result, err := h.manager.ExecuteTool(ctx, req.ToolID, req.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// getConfig 获取配置
func (h *MCPHandler) getConfig(w http.ResponseWriter, r *http.Request) {
	config := struct {
		AutoApprove bool          `json:"auto_approve"`
		Timeout     time.Duration `json:"timeout"`
	}{
		AutoApprove: h.manager.GetAutoApprove(r.Context()),
		Timeout:     h.manager.GetTimeout(r.Context()),
	}

	json.NewEncoder(w).Encode(config)
}

// updateConfig 更新配置
func (h *MCPHandler) updateConfig(w http.ResponseWriter, r *http.Request) {
	var config struct {
		AutoApprove *bool          `json:"auto_approve,omitempty"`
		Timeout     *time.Duration `json:"timeout,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config.AutoApprove != nil {
		if err := h.manager.SetAutoApprove(r.Context(), *config.AutoApprove); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if config.Timeout != nil {
		if err := h.manager.SetTimeout(r.Context(), *config.Timeout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
