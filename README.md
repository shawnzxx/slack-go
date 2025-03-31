# Slack-Go MCP Server

这是一个用 Go 语言实现的 Slack MCP (Model Context Protocol) 服务器。它提供了与 Slack API 交互的各种功能，包括发送消息、获取频道历史记录、添加表情反应等。此版本使用[mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)库实现 MCP 协议。

## 功能特性

- 列出工作区中的公共频道
- 发送消息到频道
- 在线程中回复消息
- 添加表情反应
- 获取频道历史记录
- 获取线程回复
- 获取用户列表
- 获取用户资料

## 环境要求

- Go 1.21 或更高版本
- Slack API Token
- Slack Team ID

## 环境变量

程序需要以下环境变量：

- `SLACK_TOKEN` 或 `SLACK_BOT_TOKEN`：Slack API 的访问令牌
- `SLACK_TEAM_ID`：Slack 工作区的 Team ID

## 安装

```bash
git clone https://github.com/yourusername/slack-go
cd slack-go
go mod download
```

## 编译

```bash
go build -o bin/slack-mcp ./main
```

## 运行

```bash
export SLACK_TOKEN="xoxb-your-token"
export SLACK_TEAM_ID="your-team-id"
./bin/slack-mcp
```

## 项目结构

```
slack-go/
├── main/
│   └── main.go       # 主程序入口
├── pkg/
│   └── slack/        # Slack客户端实现
│       └── client.go
├── go.mod           # Go模块定义
└── README.md        # 项目文档
```

## 技术栈

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - MCP 协议的 Go 实现
- [slack-go/slack](https://github.com/slack-go/slack) - Slack API 的 Go 客户端库

## API 工具

服务器提供以下 MCP 工具：

- `slack_list_channels`: 列出工作区中的公共频道
- `slack_post_message`: 发送消息到频道
- `slack_reply_to_thread`: 在线程中回复消息
- `slack_add_reaction`: 添加表情反应
- `slack_get_channel_history`: 获取频道历史记录
- `slack_get_thread_replies`: 获取线程回复
- `slack_get_users`: 获取用户列表
- `slack_get_user_profile`: 获取用户资料

## 许可证

MIT License
