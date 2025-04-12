# Slack-Go MCP Server

这是一个用 Go 语言实现的 Slack MCP (Model Context Protocol) 服务器。

它提供了与 Slack API 交互的各种功能，包括发送消息、获取频道历史记录、添加表情反应等。此版本使用[mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)库实现 MCP 协议。

## Tools

1. `slack_list_channels`

   - List public channels in the workspace
   - Optional inputs:
     - `limit` (number, default: 100, max: 200): Maximum number of channels to return
     - `cursor` (string): Pagination cursor for next page
   - Returns: List of channels with their IDs and information

2. `slack_post_message`

   - Post a new message to a Slack channel
   - Required inputs:
     - `channel_id` (string): The ID of the channel to post to
     - `text` (string): The message text to post
   - Returns: Message posting confirmation and timestamp

3. `slack_get_thread_replies`

   - Get all replies in a message thread
   - Required inputs:
     - `thread_url` (string): Slack message URL
       - Format: https://{workspace}.slack.com/archives/{channel_id}/{message_id}
       - URL 组成部分说明:
         - workspace: 你的 Slack 工作区名称
         - channel_id: 以 'C' 开头的频道 ID
         - message_id: 以 'p' 开头的消息 ID，包含时间戳
       - 示例:
         - 标准格式: https://myworkspace.slack.com/archives/C0123ABCDEF/p1234567890123456
         - 私有频道: https://myworkspace.slack.com/archives/C0123ABCDEF/p1234567890123456
         - 共享频道: https://myworkspace.slack.com/archives/C0123ABCDEF/p1234567890123456?thread_ts=1234567890.123456
   - Returns: List of replies with their content and metadata
   - 注意:
     - URL 可以从 Slack 客户端中通过右键点击消息并选择"Copy link"获取
     - 消息 ID 中的时间戳部分对应消息发送的 Unix 时间戳

4. `slack_get_users_profile`

   - Get detailed profile information for multiple users
   - Required inputs:
     - `user_ids` (array of strings): Array of user IDs to get profiles for
       - Format: 每个用户 ID 都以 'U' 开头，后跟数字和字母的组合
       - 示例:
         - 单个用户: ["U0123ABCDEF"]
         - 多个用户: ["U0123ABCDEF", "U9876ZYXWVU", "U5432ABCDEF"]
       - 常见错误格式:
         - ❌ 不带引号: [U0123ABCDEF]
         - ❌ 不使用数组: "U0123ABCDEF"
         - ❌ 错误前缀: ["B0123ABCDEF"] (Bot 用户使用 'B' 前缀)
         - ❌ 使用 @ 符号: ["@username"]
         - ❌ 使用邮箱: ["user@example.com"]
   - Returns: Array of user profile information including:
     - Name
     - First Name
     - Last Name
     - Real Name
     - Display Name
     - Email
     - Title
   - 注意:
     - 用户 ID 可以从 Slack 客户端中通过右键点击用户名并选择"Copy member ID"获取
     - 也可以从用户的 Slack 个人资料页面 URL 中获取
     - 每次调用最多支持 30 个用户 ID
     - 对于不存在的用户 ID 会返回错误
     - 需要确保有足够的权限访问用户资料信息
   - 使用示例:
     ```json
     {
       "user_ids": ["U0123ABCDEF", "U9876ZYXWVU"]
     }
     ```
   - 返回示例:
     ```json
     [
       {
         "name": "john.doe",
         "first_name": "John",
         "last_name": "Doe",
         "real_name": "John Doe",
         "display_name": "johndoe",
         "email": "john.doe@example.com",
         "title": "Software Engineer"
       },
       {
         "name": "jane.smith",
         "first_name": "Jane",
         "last_name": "Smith",
         "real_name": "Jane Smith",
         "display_name": "jsmith",
         "email": "jane.smith@example.com",
         "title": "Product Manager"
       }
     ]
     ```

## Environment Variables

The application requires the following environment variables:

- `SLACK_TOKEN` this token were automatically generated when you installed the app to SP Digital.
- get token from link: https://api.slack.com/apps/A08FM2YG0E5/oauth?
- `SLACK_TEAM_ID`: Your Slack workspace Team ID

### Local Testing Setup

For local testing, create a `local.env` file in the project root directory:

1. Create the file:

   ```bash
   touch local.env
   ```

2. Add your Slack credentials to `local.env`:

   ```env
   SLACK_TOKEN=xoxb-your-slack-token-here
   SLACK_TEAM_ID=your-team-id-here
   ```

   Note:

   - For `SLACK_TOKEN`, you can use either a bot token (starts with `xoxb-`) or a user token (starts with `xoxp-`)
   - The `SLACK_TEAM_ID` can be found in your Slack workspace URL or workspace settings

3. Security considerations:
   - Never commit `local.env` to version control
   - Keep your tokens secure and rotate them regularly
   - Make sure `local.env` is included in `.gitignore`

The test script (`test_single_request.sh`) will automatically load these environment variables from `local.env` when running tests.

## folder hierarchy

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

- [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - implement MCP protocol in Go
- [slack-go/slack](https://github.com/slack-go/slack) - Slack API client in Go

## license

MIT License
