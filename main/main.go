package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shawnzhang/slack-go/pkg/slack"
)

func main() {
	// set log output to stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("start slack-go MCP server...")

	// init slack client
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		token = os.Getenv("SLACK_BOT_TOKEN")
	}
	teamID := os.Getenv("SLACK_TEAM_ID")

	if token == "" || teamID == "" {
		log.Fatal("please set SLACK_TOKEN (or SLACK_BOT_TOKEN) and SLACK_TEAM_ID environment variables")
	}

	slackClient := slack.NewClient(token)

	// Create a new MCP server
	s := server.NewMCPServer(
		"slack-go",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// define tools: slack_list_channels
	listChannelsTool := mcp.NewTool("slack_list_channels",
		mcp.WithDescription("list public channels in the workspace (supports pagination)"),
		mcp.WithNumber("limit",
			mcp.Description("return the maximum number of channels (default 100, max 200)"),
			mcp.DefaultNumber(100),
		),
		mcp.WithString("cursor",
			mcp.Description("the pagination cursor for the next page results"),
		),
	)

	// define tools: slack_get_thread_replies
	getThreadRepliesTool := mcp.NewTool("slack_get_thread_replies",
		mcp.WithDescription("get all replies in a message thread"),
		mcp.WithString("thread_url",
			mcp.Required(),
			mcp.Description("Slack message URL"),
		),
	)

	// define tools: postMessageTool
	postMessageTool := mcp.NewTool("post_message",
		mcp.WithDescription("post a message to a Slack channel"),
		mcp.WithString("channel_id",
			mcp.Required(),
			mcp.Description("ID of the channel to post the message to"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text of the message to post"),
		),
	)

	// define tools: slack_get_users_profile
	getUsersProfileTool := mcp.NewTool("slack_get_users_profile",
		mcp.WithDescription("get multiple users' profile information"),
		mcp.WithArray("user_ids",
			mcp.Required(),
			mcp.Description("Array of user IDs to get profiles for"),
		),
	)

	// add tools and handle functions
	s.AddTool(listChannelsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := 100
		if l, ok := request.Params.Arguments["limit"].(float64); ok {
			limit = int(l)
			log.Printf("use provided limit: %d", limit)
		} else {
			log.Printf("use default limit: %d", limit)
		}

		cursor := ""
		if c, ok := request.Params.Arguments["cursor"].(string); ok {
			cursor = c
			log.Printf("use provided cursor: %s", cursor)
		}

		log.Printf("start to get channel list...")

		// call slack api to get channel list
		result, err := slackClient.ListChannels(limit, cursor)
		if err != nil {
			log.Printf("failed to get channel list: %v", err)
			return nil, fmt.Errorf("failed to get channel list: %v", err)
		}
		log.Printf("success to get channel list")

		channelsJSON, err := json.Marshal(result.Channels)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize channel list: %v", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("channel list: \n%s", string(channelsJSON))), nil
	})

	s.AddTool(getThreadRepliesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		threadURL, ok := request.Params.Arguments["thread_url"].(string)
		if !ok || threadURL == "" {
			log.Printf("error: invalid thread_url: %v", request.Params.Arguments["thread_url"])
			return nil, fmt.Errorf("thread_url is required")
		}
		log.Printf("process thread URL: %s", threadURL)

		// call slack api to get thread replies
		log.Printf("start to get thread replies...")
		result, err := slackClient.GetThreadReplies(threadURL)
		if err != nil {
			log.Printf("failed to get thread replies: %v", err)
			return nil, fmt.Errorf("failed to get thread replies: %v", err)
		}
		log.Printf("success to get thread replies")

		messagesJSON, err := json.Marshal(result.Messages)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize thread replies: %v", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("thread replies: \n%s", string(messagesJSON))), nil
	})

	s.AddTool(postMessageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, ok := request.Params.Arguments["channel_id"].(string)
		if !ok || channelID == "" {
			log.Printf("error: invalid channel_id: %v", request.Params.Arguments["channel_id"])
			return nil, fmt.Errorf("channel_id is required")
		}

		text, ok := request.Params.Arguments["text"].(string)
		if !ok || text == "" {
			log.Printf("error: invalid text: %v", request.Params.Arguments["text"])
			return nil, fmt.Errorf("text is required")
		}

		log.Printf("posting message to channel: %s", channelID)

		// call slack api to post message
		message, err := slackClient.PostMessage(channelID, text)
		if err != nil {
			log.Printf("failed to post message: %v", err)
			return nil, fmt.Errorf("failed to post message: %v", err)
		}
		log.Printf("success to post message")

		messageJSON, err := json.Marshal(message)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize message: %v", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("message posted: \n%s", string(messageJSON))), nil
	})

	s.AddTool(getUsersProfileTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取并验证用户ID数组
		userIDsInterface, ok := request.Params.Arguments["user_ids"].([]interface{})
		if !ok || len(userIDsInterface) == 0 {
			log.Printf("error: invalid user_ids: %v", request.Params.Arguments["user_ids"])
			return nil, fmt.Errorf("user_ids array is required and cannot be empty")
		}

		// 转换interface{}数组为string数组
		userIDs := make([]string, len(userIDsInterface))
		for i, v := range userIDsInterface {
			userID, ok := v.(string)
			if !ok || userID == "" {
				return nil, fmt.Errorf("invalid user ID at position %d", i)
			}
			userIDs[i] = userID
		}

		log.Printf("getting profiles for users: %v", userIDs)

		// 调用slack api获取多个用户的资料
		profiles, err := slackClient.GetFilteredUsersProfile(userIDs)
		if err != nil {
			log.Printf("failed to get user profiles: %v", err)
			return nil, fmt.Errorf("failed to get user profiles: %v", err)
		}
		log.Printf("success to get user profiles")

		// 序列化结果
		profilesJSON, err := json.Marshal(profiles)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize user profiles: %v", err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("user profiles: \n%s", string(profilesJSON))), nil
	})

	// start standard input/output server
	log.Printf("MCP server is ready, start to process requests...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
