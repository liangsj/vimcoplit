package api

import (
	"encoding/json"
	"net/http"

	"github.com/liangsj/vimcoplit/internal/core"
	"github.com/liangsj/vimcoplit/internal/models"
)

// Handler 处理所有HTTP请求
type Handler struct {
	service core.Service
}

// NewHandler 创建新的API处理器
func NewHandler(service core.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ServeHTTP 实现http.Handler接口
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 路由处理
	switch r.URL.Path {
	case "/api/tasks":
		h.handleTasks(w, r)
	case "/api/files":
		h.handleFiles(w, r)
	case "/api/execute":
		h.handleExecute(w, r)
	case "/api/generate":
		h.handleGenerate(w, r)
	case "/api/model":
		h.handleModel(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleTasks 处理任务相关的请求
func (h *Handler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var req struct {
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		taskID, err := h.service.CreateTask(r.Context(), req.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"task_id": taskID})

	case "GET":
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, "task ID is required", http.StatusBadRequest)
			return
		}

		task, err := h.service.GetTask(r.Context(), taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleFiles 处理文件操作相关的请求
func (h *Handler) handleFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "path is required", http.StatusBadRequest)
			return
		}

		content, err := h.service.ReadFile(r.Context(), path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"content": content})

	case "POST":
		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.service.WriteFile(r.Context(), req.Path, req.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "path is required", http.StatusBadRequest)
			return
		}

		if err := h.service.DeleteFile(r.Context(), path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleExecute 处理命令执行请求
func (h *Handler) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.service.ExecuteCommand(r.Context(), req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"output": output})
}

// handleGenerate 处理AI响应生成请求
func (h *Handler) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TaskID string `json:"task_id"`
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.GenerateResponse(r.Context(), req.TaskID, req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

// handleModel 处理模型相关的请求
func (h *Handler) handleModel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// 获取当前模型
		modelType := h.service.GetCurrentModel()
		json.NewEncoder(w).Encode(map[string]string{"model": string(modelType)})

	case "POST":
		// 切换模型
		var req struct {
			ModelType string `json:"model_type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.service.SwitchModel(r.Context(), models.ModelType(req.ModelType)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
