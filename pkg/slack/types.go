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

// UserProfileInfo represents filtered user profile information
type UserProfileInfo struct {
	// Name is the username of the Slack user (e.g. johndoe)
	Name string `json:"name"`
	// FullName is the actual name of the Slack user (e.g. John Doe)
	FullName string `json:"full_name"`
	// DisplayName is the display name set by the user in their profile
	DisplayName string `json:"display_name"`
	// Email is the user's email address
	Email string `json:"email"`
	// Title is the user's job title or role in the organization
	Title string `json:"title"`
}

// GetThreadRepliesResponse represents the response from a GetThreadReplies call
type GetThreadRepliesResponse struct {
	Messages []slack.Message
}
