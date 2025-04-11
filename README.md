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

## project structure

### folder hierarchy

```
slack-go/
├── main/
│ └── main.go # Main entry point of the application
├── pkg/
│ └── slack/ # Implementation of the Slack client
│ └── client.go
├── vendor/ # Vendor directory for dependencies
├── go.mod # Go module definition
├── go.sum # Go module dependencies checksum
├── Makefile # Makefile for build automation
├── README.md # Project documentation
└── test_single_request.sh # Script for testing single requests to the Slack MCP server
```

### folder description

- **main/**: Contains the main entry point of the application where the server is initialized and started.
- **pkg/slack/**: Contains the implementation of the Slack client, which wraps the Slack API functionalities.
- **vendor/**: Holds the vendored dependencies to ensure consistent builds.
- **go.mod**: Defines the module's dependencies and versions.
- **go.sum**: Contains checksums for the module's dependencies.
- **Makefile**: Provides build automation tasks such as building the binary and managing dependencies.
- **README.md**: Provides documentation about the project, including setup, usage, and features.
- **test_single_request.sh**: A script to test various requests to the Slack MCP server, ensuring the server's functionalities are working as expected.

## tech stack

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - 实现 MCP 协议的 Go 库
- [slack-go/slack](https://github.com/slack-go/slack) - Slack API 的 Go 客户端库

## API tools

The server provides the following MCP tools:

- `slack_list_channels`: List public channels in the workspace
- `slack_post_message`: Send a message to a channel
- `slack_reply_to_thread`: Reply to a message in a thread
- `slack_add_reaction`: Add an emoji reaction
- `slack_get_channel_history`: Get the history of a channel
- `slack_get_thread_replies`: Get replies to a message in a thread
- `slack_get_users`: Get the list of users
- `slack_get_user_profile`: Get the profile of a user

## license

MIT License
