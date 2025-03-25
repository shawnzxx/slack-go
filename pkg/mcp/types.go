package mcp

import "encoding/json"

// JSONRPCRequest 表示一个JSON-RPC请求
type JSONRPCRequest struct {
	ID      interface{} `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse 表示一个JSON-RPC响应
type JSONRPCResponse struct {
	ID      interface{}   `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError 表示JSON-RPC错误
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ServerInfo 表示服务器信息
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Capabilities 表示服务器能力
type Capabilities struct {
	Tools map[string]interface{} `json:"tools"`
}

// InitializeResult 表示初始化结果
type InitializeResult struct {
	Capabilities    Capabilities `json:"capabilities"`
	ProtocolVersion string       `json:"protocolVersion"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

// Tool 表示一个工具
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// Types for resources
type ListResourcesResult struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// Types for prompts
type ListPromptsResult struct {
	Prompts []Prompt `json:"prompts"`
}

type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// ListToolsResult 表示工具列表结果
type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

// ToolContent 表示工具调用的响应内容
type ToolContent struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// CallToolResult 表示工具调用的结果
type CallToolResult struct {
	Content []ToolContent `json:"content"`
}

// 错误码常量
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)
