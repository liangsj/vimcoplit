# VimCoplit

VimCoplit 是一个基于 [Cline](https://github.com/cline/cline) 的 Neovim 实现版本，它是一个强大的 AI 编码助手，直接集成到您的编辑器中。本项目旨在为 Neovim 用户带来 Cline 的强大功能，并使用 Go 后端提供高性能和可靠性。

## 功能特点

- **AI 驱动的编码助手**：支持多种大语言模型
  - Claude 3.5 Sonnet：强大的代码理解和生成能力
  - 豆包：专注于中文场景的智能助手
  - DeepSeek：高性能的开源模型
- **文件操作**：创建、编辑和管理文件，获得 AI 辅助
- **终端集成**：执行命令并监控输出
- **浏览器自动化**：使用无头浏览器测试 Web 应用
- **自定义工具**：通过模型上下文协议（MCP）扩展功能
- **上下文管理**：添加 URL、问题、文件和文件夹以提供 AI 上下文

## 系统架构

VimCoplit 由两个主要组件组成：

1. **Neovim 插件**：使用 Lua 编写，提供用户界面和 Neovim 集成
2. **Go 后端**：处理 AI 交互、文件操作和系统命令

## 系统要求

- Neovim 0.9.0 或更高版本
- Go 1.21 或更高版本
- 支持的 AI 模型 API 密钥：
  - [Claude API 密钥](https://console.anthropic.com/)
  - [豆包 API 密钥](https://www.doubao.com/)
  - [DeepSeek API 密钥](https://platform.deepseek.com/)

## 安装

```bash
# 使用您喜欢的 Neovim 插件管理器
# 例如，使用 lazy.nvim：
{
  "liangsj/vimcoplit",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim",
  },
  build = "go build -o bin/vimcoplit ./cmd/vimcoplit",
}
```

## 配置

```lua
-- 在您的 Neovim 配置中
require('vimcoplit').setup({
  -- 选择使用的模型
  model = "claude-3-sonnet-20240229", -- 可选: "claude-3-sonnet-20240229", "doubao", "deepseek"
  
  -- 模型 API 密钥
  api_key = "your-api-key",
  
  -- 其他配置选项
  max_tokens = 4096,
  temperature = 0.7,
})
```

## 使用方法

- `:VimCoplit` - 打开 VimCoplit 界面
- `:VimCoplitTask <task>` - 开始新任务
- `:VimCoplitAddContext` - 为当前任务添加上下文
- `:VimCoplitSwitchModel` - 切换使用的 AI 模型

## 开发

### 项目结构

```
vimcoplit/
├── cmd/                # Go 后端入口点
├── internal/           # Go 后端实现
│   ├── api/           # API 处理器
│   ├── core/          # 核心服务
│   └── models/        # AI 模型集成
├── lua/                # Neovim 插件代码
│   └── vimcoplit/      # 插件模块
├── scripts/            # 构建和工具脚本
└── test/               # 测试文件
```

### 构建

```bash
# 构建 Go 后端
go build -o bin/vimcoplit ./cmd/vimcoplit

# 构建 Neovim 插件
nvim --headless -c "luafile scripts/build.lua" -c "quit"
```

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

Apache 2.0 © 2025 VimCoplit Contributors

## 致谢

本项目基于 [Cline](https://github.com/cline/cline) 项目开发，感谢原项目的所有贡献者。 