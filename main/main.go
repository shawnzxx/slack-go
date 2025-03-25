package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/shawnzhang/slack-go/pkg/mcp"
	"github.com/shawnzhang/slack-go/pkg/slack"
)

func main() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	// 设置日志输出到stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("启动 slack-go MCP 服务器...")

	// 初始化Slack客户端
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		token = os.Getenv("SLACK_BOT_TOKEN")
	}
	teamID := os.Getenv("SLACK_TEAM_ID")

	if token == "" || teamID == "" {
		log.Fatal("请设置 SLACK_TOKEN (或 SLACK_BOT_TOKEN) 和 SLACK_TEAM_ID 环境变量")
	}

	slackClient := slack.NewClient(token)

	for {
		var request mcp.JSONRPCRequest
		// 从标准输入读取请求
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				log.Printf("收到EOF，服务器正常退出")
				os.Exit(0)
			}
			log.Printf("解码请求错误: %v", err)
			sendError(encoder, nil, mcp.ParseError, "解析 JSON 失败")
			continue
		}

		log.Printf("收到请求: %v", PrettyJSON(request))

		if request.JSONRPC != "2.0" {
			sendError(encoder, request.ID, mcp.InvalidRequest, "仅支持 JSON-RPC 2.0")
			continue
		}

		var response interface{}

		switch request.Method {
		case "initialize":
			response = mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Result: mcp.InitializeResult{
					ProtocolVersion: "2024-11-05",
					ServerInfo: mcp.ServerInfo{
						Name:    "slack-go",
						Version: "1.0.0",
					},
					Capabilities: mcp.Capabilities{
						Tools: map[string]interface{}{},
					},
				},
			}

		case "notifications/initialized", "initialized":
			log.Printf("服务器初始化成功")
			continue // 跳过发送响应，因为通知不需要响应

		case "tools/list":
			response = mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Result: mcp.ListToolsResult{
					Tools: []mcp.Tool{
						{
							Name:        "slack_list_channels",
							Description: "列出工作区中的公共频道（支持分页）",
							InputSchema: json.RawMessage(`{
								"type": "object",
								"properties": {
									"limit": {
										"type": "number",
										"description": "返回的最大频道数量（默认100，最大200）",
										"default": 100
									},
									"cursor": {
										"type": "string",
										"description": "用于下一页结果的分页游标"
									}
								}
							}`),
						},
						{
							Name:        "slack_get_thread_replies",
							Description: "获取消息线程中的所有回复",
							InputSchema: json.RawMessage(`{
								"type": "object",
								"properties": {
									"thread_url": {
										"type": "string",
										"description": "Slack消息URL"
									}
								},
								"required": ["thread_url"]
							}`),
						},
					},
				},
			}

		case "tools/call":
			log.Printf("处理 tools/call 请求")
			params, ok := request.Params.(map[string]interface{})
			if !ok {
				log.Printf("错误: 无效的参数类型: %T", request.Params)
				sendError(encoder, request.ID, mcp.InvalidParams, "无效的参数")
				continue
			}

			toolName, ok := params["name"].(string)
			if !ok {
				log.Printf("错误: 工具名称未找到或类型无效: %T", params["name"])
				sendError(encoder, request.ID, mcp.InvalidParams, "无效的工具名称")
				continue
			}
			log.Printf("请求的工具: %s", toolName)

			args, ok := params["arguments"].(map[string]interface{})
			if !ok {
				log.Printf("错误: 无效的参数类型: %T", params["arguments"])
				sendError(encoder, request.ID, mcp.InvalidParams, "无效的参数")
				continue
			}
			log.Printf("收到参数: %v", PrettyJSON(args))

			switch toolName {
			case "slack_list_channels":
				limit := 100
				if l, ok := args["limit"].(float64); ok {
					limit = int(l)
					log.Printf("使用提供的限制: %d", limit)
				} else {
					log.Printf("使用默认限制: %d", limit)
				}

				cursor := ""
				if c, ok := args["cursor"].(string); ok {
					cursor = c
					log.Printf("使用提供的游标: %s", cursor)
				}

				log.Printf("开始获取频道列表...")
				result, err := slackClient.ListChannels(limit, cursor)
				if err != nil {
					log.Printf("获取频道列表失败: %v", err)
					sendError(encoder, request.ID, mcp.InternalError, fmt.Sprintf("获取频道列表失败: %v", err))
					continue
				}
				log.Printf("成功获取频道列表")

				response = mcp.JSONRPCResponse{
					JSONRPC: "2.0",
					ID:      request.ID,
					Result: mcp.CallToolResult{
						Content: []mcp.ToolContent{
							{
								Type:     "text",
								Text:     fmt.Sprintf("频道列表：\n%s", string(mustMarshalJSON(result.Channels))),
								MimeType: "text/plain",
							},
						},
					},
				}

			case "slack_get_thread_replies":
				threadURL, ok := args["thread_url"].(string)
				if !ok || threadURL == "" {
					log.Printf("错误: 无效的thread_url: %v", args["thread_url"])
					sendError(encoder, request.ID, mcp.InvalidParams, "thread_url是必需的")
					continue
				}
				log.Printf("处理线程URL: %s", threadURL)

				log.Printf("开始获取线程回复...")
				result, err := slackClient.GetThreadReplies(threadURL)
				if err != nil {
					log.Printf("获取线程回复失败: %v", err)
					sendError(encoder, request.ID, mcp.InternalError, fmt.Sprintf("获取线程回复失败: %v", err))
					continue
				}
				log.Printf("成功获取线程回复")

				response = mcp.JSONRPCResponse{
					JSONRPC: "2.0",
					ID:      request.ID,
					Result: mcp.CallToolResult{
						Content: []mcp.ToolContent{
							{
								Type:     "text",
								Text:     fmt.Sprintf("线程回复：\n%s", string(mustMarshalJSON(result.Messages))),
								MimeType: "text/plain",
							},
						},
					},
				}

			default:
				log.Printf("错误: 未知的工具: %s", toolName)
				sendError(encoder, request.ID, mcp.MethodNotFound, "未知的工具")
				continue
			}

		case "cancelled":
			if params, ok := request.Params.(map[string]interface{}); ok {
				log.Printf("收到取消通知，请求ID: %v, 原因: %v",
					params["requestId"], params["reason"])
			} else {
				log.Printf("收到取消通知，参数无效")
			}
			continue // 跳过发送响应，因为通知不需要响应

		default:
			sendError(encoder, request.ID, mcp.MethodNotFound, "方法未实现")
			continue
		}

		// 发送响应
		log.Printf("发送响应: %v", PrettyJSON(response))
		sendResponse(encoder, response)
	}

	log.Printf("slack-go MCP 服务器退出循环...")
}

// 发送错误响应
func sendError(encoder *json.Encoder, id interface{}, code int, message string) {
	// 对于请求中的null ID，我们应该响应null ID
	var responseID int
	if id != nil {
		if idFloat, ok := id.(float64); ok {
			responseID = int(idFloat)
		}
	}

	response := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      responseID,
		Error: &mcp.JSONRPCError{
			Code:    code,
			Message: message,
		},
	}

	log.Printf("发送错误响应: %v", PrettyJSON(response))
	if err := encoder.Encode(response); err != nil {
		log.Printf("编码错误响应失败: %v", err)
	}
}

// 发送响应
func sendResponse(encoder *json.Encoder, response interface{}) {
	if err := encoder.Encode(response); err != nil {
		log.Printf("编码响应失败: %v", err)
	}
}

// 辅助函数：将数据转换为格式化的JSON字符串
func PrettyJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("错误: %v", err)
	}
	return string(b)
}

// 辅助函数：将数据转换为JSON字符串
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("JSON编码失败: %v", err)
		return []byte("{}")
	}
	return data
}
