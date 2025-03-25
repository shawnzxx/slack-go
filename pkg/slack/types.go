package slack

import (
	"github.com/slack-go/slack"
)

// Message represents a Slack message
type Message struct {
	Timestamp       string
	Channel         string
	Text            string
	ThreadTimestamp string
}

// GetUsersResponse represents the response from a GetUsers call
type GetUsersResponse struct {
	Members          []slack.User
	ResponseMetadata struct {
		NextCursor string
	}
}

// GetUsersParameters represents the parameters for a GetUsers call
type GetUsersParameters struct {
	Limit  int
	Cursor string
}

// GetConversationsResponse represents the response from a GetConversations call
type GetConversationsResponse struct {
	Channels         []slack.Channel
	ResponseMetadata struct {
		NextCursor string
	}
}

// GetThreadRepliesResponse represents the response from a GetThreadReplies call
type GetThreadRepliesResponse struct {
	Messages []slack.Message
}
