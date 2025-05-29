package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPExecutor 是一个基于 HTTP 的工具执行器
type HTTPExecutor struct {
	client *http.Client
}

// NewHTTPExecutor 创建一个新的 HTTP 执行器
func NewHTTPExecutor(timeout time.Duration) *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Execute 执行工具
func (e *HTTPExecutor) Execute(ctx context.Context, tool *Tool, params map[string]interface{}) (*ToolExecutionResult, error) {
	// 验证参数
	if err := tool.ValidateParameters(params); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %v", err)
	}

	// 准备请求
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %v", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", tool.Metadata["endpoint"], bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if auth := tool.Metadata["auth"]; auth != "" {
		req.Header.Set("Authorization", auth)
	}

	// 执行请求
	startTime := time.Now()
	resp, err := e.client.Do(req)
	endTime := time.Now()

	result := &ToolExecutionResult{
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err != nil {
		if err == context.DeadlineExceeded {
			result.Status = ToolExecutionStatusTimeout
			result.Error = "execution timed out"
		} else {
			result.Status = ToolExecutionStatusError
			result.Error = fmt.Sprintf("request failed: %v", err)
		}
		return result, nil
	}
	defer resp.Body.Close()

	// 解析响应
	var responseBody interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		result.Status = ToolExecutionStatusError
		result.Error = fmt.Sprintf("failed to decode response: %v", err)
		return result, nil
	}

	// 检查响应状态
	if resp.StatusCode >= 400 {
		result.Status = ToolExecutionStatusError
		if errMsg, ok := responseBody.(map[string]interface{})["error"]; ok {
			result.Error = fmt.Sprintf("%v", errMsg)
		} else {
			result.Error = fmt.Sprintf("server returned status %d", resp.StatusCode)
		}
		return result, nil
	}

	// 成功
	result.Status = ToolExecutionStatusSuccess
	result.Result = responseBody
	return result, nil
}

// LocalExecutor 是一个本地工具执行器
type LocalExecutor struct {
	handlers map[string]ToolHandler
}

// ToolHandler 定义了本地工具处理函数
type ToolHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// NewLocalExecutor 创建一个新的本地执行器
func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{
		handlers: make(map[string]ToolHandler),
	}
}

// RegisterHandler 注册工具处理函数
func (e *LocalExecutor) RegisterHandler(toolID string, handler ToolHandler) {
	e.handlers[toolID] = handler
}

// Execute 执行工具
func (e *LocalExecutor) Execute(ctx context.Context, tool *Tool, params map[string]interface{}) (*ToolExecutionResult, error) {
	// 验证参数
	if err := tool.ValidateParameters(params); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %v", err)
	}

	// 查找处理函数
	handler, exists := e.handlers[tool.ID]
	if !exists {
		return nil, fmt.Errorf("no handler registered for tool %s", tool.ID)
	}

	// 执行处理函数
	startTime := time.Now()
	result, err := handler(ctx, params)
	endTime := time.Now()

	execResult := &ToolExecutionResult{
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err != nil {
		if err == context.DeadlineExceeded {
			execResult.Status = ToolExecutionStatusTimeout
			execResult.Error = "execution timed out"
		} else {
			execResult.Status = ToolExecutionStatusError
			execResult.Error = err.Error()
		}
		return execResult, nil
	}

	execResult.Status = ToolExecutionStatusSuccess
	execResult.Result = result
	return execResult, nil
}
