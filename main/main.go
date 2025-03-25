package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shawnzhang/slack-go/pkg/mcp"
	"github.com/shawnzhang/slack-go/pkg/slack"
)

// 设置SSE响应头
func setSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

// 错误处理函数
func handleError(w http.ResponseWriter, status int, message string) {
	setSSEHeaders(w)
	w.WriteHeader(status)
	writeSSEEvent(w, "error", map[string]string{
		"message": message,
	})
}

func main() {
	// 设置日志输出到stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("Starting slack-go MCP server...")

	// 初始化Slack客户端
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		token = os.Getenv("SLACK_BOT_TOKEN")
	}
	teamID := os.Getenv("SLACK_TEAM_ID")

	if token == "" || teamID == "" {
		log.Fatal("Please set SLACK_TOKEN (or SLACK_BOT_TOKEN) and SLACK_TEAM_ID environment variables")
	}

	slackClient := slack.NewClient(token)

	// 创建路由器
	mux := http.NewServeMux()

	// 添加CORS中间件
	corsHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 设置CORS头
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// 对于SSE请求，先设置SSE头
			if r.Header.Get("Accept") == "text/event-stream" {
				setSSEHeaders(w)
			}

			next(w, r)
		}
	}

	// 健康检查端点
	mux.HandleFunc("/health", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	}))

	// MCP端点
	mux.HandleFunc("/mcp/initialize", corsHandler(handleInitialize))
	mux.HandleFunc("/mcp/tools/list", corsHandler(handleToolsList))
	mux.HandleFunc("/mcp/slack/list-channels", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		handleSlackListChannels(w, r, slackClient)
	}))
	mux.HandleFunc("/mcp/slack/get-thread-replies", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		handleSlackGetThreadReplies(w, r, slackClient)
	}))

	// 创建服务器
	server := &http.Server{
		Addr:              ":3333",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 启动HTTP服务器
	go func() {
		log.Printf("Starting HTTP server on %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 优雅关闭服务器
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

// SSE工具函数
func writeSSEEvent(w http.ResponseWriter, event string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling SSE data: %v", err)
		fmt.Fprintf(w, "event: error\ndata: {\"message\": \"Internal server error\"}\n\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	// 发送心跳注释
	fmt.Fprintf(w, ": heartbeat\n\n")

	// 发送实际事件
	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)

	// 确保数据被发送
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// 处理初始化请求
func handleInitialize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	result := mcp.InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo: mcp.ServerInfo{
			Name:    "slack-go",
			Version: "1.0.0",
		},
		Capabilities: mcp.Capabilities{
			Tools: map[string]interface{}{},
		},
	}

	writeSSEEvent(w, "initialize", result)
}

// 处理工具列表请求
func handleToolsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tools := []mcp.Tool{
		{
			Name:        "slack_list_channels",
			Description: "List public channels in the workspace with pagination",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"limit": map[string]interface{}{
						"type":        "number",
						"description": "Maximum number of channels to return (default 100, max 200)",
						"default":     100,
					},
					"cursor": map[string]interface{}{
						"type":        "string",
						"description": "Pagination cursor for next page of results",
					},
				},
			},
		},
		{
			Name:        "slack_get_thread_replies",
			Description: "Get all replies in a message thread",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"thread_url": map[string]interface{}{
						"type":        "string",
						"description": "The Slack message URL",
					},
				},
				"required": []string{"thread_url"},
			},
		},
	}

	writeSSEEvent(w, "tools/list", mcp.ListToolsResult{Tools: tools})
}

// 处理Slack频道列表请求
func handleSlackListChannels(w http.ResponseWriter, r *http.Request, client *slack.Client) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从查询参数获取limit和cursor
	query := r.URL.Query()
	limit := 100
	if l := query.Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	cursor := query.Get("cursor")

	// 发送进度事件
	writeSSEEvent(w, "progress", map[string]string{
		"status": "开始获取频道列表...",
	})

	// 获取频道列表
	result, err := client.ListChannels(limit, cursor)
	if err != nil {
		writeSSEEvent(w, "error", map[string]string{
			"message": fmt.Sprintf("获取频道列表失败: %v", err),
		})
		return
	}

	writeSSEEvent(w, "result", result)
	writeSSEEvent(w, "done", nil)
}

// 处理Slack线程回复请求
func handleSlackGetThreadReplies(w http.ResponseWriter, r *http.Request, client *slack.Client) {
	if r.Method != http.MethodGet {
		handleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	threadURL := r.URL.Query().Get("thread_url")
	if threadURL == "" {
		writeSSEEvent(w, "error", map[string]string{
			"message": "缺少thread_url参数",
		})
		return
	}

	// 发送进度事件
	writeSSEEvent(w, "progress", map[string]string{
		"status": "开始获取线程回复...",
	})

	// 获取线程回复
	result, err := client.GetThreadReplies(threadURL)
	if err != nil {
		writeSSEEvent(w, "error", map[string]string{
			"message": fmt.Sprintf("获取线程回复失败: %v", err),
		})
		return
	}

	writeSSEEvent(w, "result", result)
	writeSSEEvent(w, "done", nil)
}
