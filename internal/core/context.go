package core

import (
	"errors"
	"sync"
	"time"
)

// ContextType 表示上下文的类型
type ContextType string

const (
	ContextTypeURL      ContextType = "url"
	ContextTypeQuestion ContextType = "question"
	ContextTypeFile     ContextType = "file"
	ContextTypeFolder   ContextType = "folder"
)

// ContextItem 表示一个上下文条目
type ContextItem interface {
	GetID() string
	GetType() ContextType
	GetValue() string
	GetCreatedAt() time.Time
}

// BaseContextItem 提供通用字段
type BaseContextItem struct {
	ID        string
	Type      ContextType
	Value     string
	CreatedAt time.Time
}

func (b *BaseContextItem) GetID() string           { return b.ID }
func (b *BaseContextItem) GetType() ContextType    { return b.Type }
func (b *BaseContextItem) GetValue() string        { return b.Value }
func (b *BaseContextItem) GetCreatedAt() time.Time { return b.CreatedAt }

// NewContextItem 创建一个新的上下文条目
func NewContextItem(id string, typ ContextType, value string) ContextItem {
	return &BaseContextItem{
		ID:        id,
		Type:      typ,
		Value:     value,
		CreatedAt: time.Now(),
	}
}

// ContextManager 定义了上下文管理器的接口
type ContextManager interface {
	AddItem(item ContextItem)
	RemoveItem(id string) error
	GetItem(id string) (ContextItem, error)
	ListItems() []ContextItem
}

// Manager 是 ContextManager 接口的具体实现
type Manager struct {
	mu    sync.RWMutex
	items map[string]ContextItem // key: id
}

// NewManager 创建一个新的上下文管理器
func NewManager() ContextManager {
	return &Manager{
		items: make(map[string]ContextItem),
	}
}

// AddItem 添加一个上下文项
func (m *Manager) AddItem(item ContextItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[item.GetID()] = item
}

// RemoveItem 删除一个上下文项
func (m *Manager) RemoveItem(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[id]; !ok {
		return errors.New("context item not found")
	}
	delete(m.items, id)
	return nil
}

// GetItem 查询一个上下文项
func (m *Manager) GetItem(id string) (ContextItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("context item not found")
	}
	return item, nil
}

// ListItems 列出所有上下文项
func (m *Manager) ListItems() []ContextItem {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]ContextItem, 0, len(m.items))
	for _, item := range m.items {
		result = append(result, item)
	}
	return result
}
