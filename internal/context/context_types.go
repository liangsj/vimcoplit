package context

import (
	"time"

	"github.com/liangsj/vimcoplit/internal/core"
)

// ContextType 表示上下文的类型
// 支持 URL、问题、文件、文件夹
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
	Type      core.ContextType
	Value     string
	CreatedAt time.Time
}

func (b *BaseContextItem) GetID() string             { return b.ID }
func (b *BaseContextItem) GetType() core.ContextType { return b.Type }
func (b *BaseContextItem) GetValue() string          { return b.Value }
func (b *BaseContextItem) GetCreatedAt() time.Time   { return b.CreatedAt }

// NewContextItem 创建一个新的上下文条目
func NewContextItem(id string, typ core.ContextType, value string) core.ContextItem {
	return &BaseContextItem{
		ID:        id,
		Type:      typ,
		Value:     value,
		CreatedAt: time.Now(),
	}
}
