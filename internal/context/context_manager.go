package context

import (
	"errors"
	"sync"

	"github.com/liangsj/vimcoplit/internal/core"
)

// Manager 是 ContextManager 接口的具体实现
type Manager struct {
	mu    sync.RWMutex
	items map[string]core.ContextItem // key: id
}

// NewManager 创建一个新的上下文管理器
func NewManager() core.ContextManager {
	return &Manager{
		items: make(map[string]core.ContextItem),
	}
}

// AddItem 添加一个上下文项
func (m *Manager) AddItem(item core.ContextItem) {
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
func (m *Manager) GetItem(id string) (core.ContextItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("context item not found")
	}
	return item, nil
}

// ListItems 列出所有上下文项
func (m *Manager) ListItems() []core.ContextItem {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]core.ContextItem, 0, len(m.items))
	for _, item := range m.items {
		result = append(result, item)
	}
	return result
}
