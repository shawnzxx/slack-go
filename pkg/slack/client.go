package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

// Client wraps the slack client with our custom methods
type Client struct {
	api *slack.Client
}

// NewClient creates a new Slack client
func NewClient(token string) *Client {
	return &Client{
		api: slack.New(token),
	}
}

// PostMessage posts a message to a channel
func (c *Client) PostMessage(channelID, text string) (*Message, error) {
	_, timestamp, err := c.api.PostMessage(
		channelID,
		slack.MsgOptionText(text, false),
	)
	if err != nil {
		return nil, err
	}

	return &Message{
		Timestamp: timestamp,
		Channel:   channelID,
		Text:      text,
	}, nil
}

// PostReply posts a reply to a thread
func (c *Client) PostReply(channelID, threadTS, text string) (*Message, error) {
	_, timestamp, err := c.api.PostMessage(
		channelID,
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(threadTS),
	)
	if err != nil {
		return nil, err
	}

	return &Message{
		Timestamp:       timestamp,
		Channel:         channelID,
		Text:            text,
		ThreadTimestamp: threadTS,
	}, nil
}

// AddReaction adds a reaction to a message
func (c *Client) AddReaction(channelID, timestamp, reaction string) error {
	return c.api.AddReaction(reaction, slack.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	})
}

// GetChannelHistory gets the message history of a channel
func (c *Client) GetChannelHistory(channelID string, limit int) (*slack.GetConversationHistoryResponse, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
	}
	return c.api.GetConversationHistory(params)
}

// GetThreadReplies gets all replies in a thread
func (c *Client) GetThreadReplies(threadURL string) (*GetThreadRepliesResponse, error) {
	// 解析URL获取channelID和timestamp
	// URL格式: https://workspace.slack.com/archives/C0734812MFG/p1742788004223029
	parts := strings.Split(threadURL, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid thread URL format")
	}

	channelID := parts[len(parts)-2]
	timestampStr := parts[len(parts)-1]

	// 处理timestamp格式
	// 将格式从 p1742788004223029 转换为 1742788004.223029
	if !strings.HasPrefix(timestampStr, "p") {
		return nil, fmt.Errorf("invalid timestamp format in URL")
	}
	tsNum := timestampStr[1:] // 去掉p前缀
	if len(tsNum) != 16 {
		return nil, fmt.Errorf("invalid timestamp length")
	}
	threadTS := fmt.Sprintf("%s.%s", tsNum[:10], tsNum[10:])

	params := &slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: threadTS,
	}
	messages, _, _, err := c.api.GetConversationReplies(params)
	if err != nil {
		return nil, err
	}
	return &GetThreadRepliesResponse{
		Messages: messages,
	}, nil
}

// GetUsers gets a list of all users
func (c *Client) GetUsers(limit int, cursor string) (*GetUsersResponse, error) {
	users, err := c.api.GetUsers()
	if err != nil {
		return nil, err
	}
	return &GetUsersResponse{
		Members: users,
		ResponseMetadata: struct{ NextCursor string }{
			NextCursor: cursor,
		},
	}, nil
}

// GetUserProfile gets a user's profile
func (c *Client) GetUserProfile(userID string) (*slack.UserProfile, error) {
	user, err := c.api.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}
	return &user.Profile, nil
}

// GetFilteredUserProfile gets filtered user profile information
func (c *Client) GetFilteredUserProfile(userID string) (*UserProfileInfo, error) {
	user, err := c.api.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	return &UserProfileInfo{
		Name:        user.Name,
		FullName:    user.Profile.RealName,
		DisplayName: user.Profile.DisplayName,
		Email:       user.Profile.Email,
		Title:       user.Profile.Title,
	}, nil
}

// GetFilteredUsersProfile gets filtered user profile information for multiple users
func (c *Client) GetFilteredUsersProfile(userIDs []string) ([]*UserProfileInfo, error) {
	users, err := c.api.GetUsersInfo(userIDs...)
	if err != nil {
		return nil, err
	}

	profiles := make([]*UserProfileInfo, 0, len(*users))
	for _, user := range *users {
		profiles = append(profiles, &UserProfileInfo{
			Name:        user.Name,
			FullName:    user.Profile.RealName,
			DisplayName: user.Profile.DisplayName,
			Email:       user.Profile.Email,
			Title:       user.Profile.Title,
		})
	}

	return profiles, nil
}

// ListChannels lists all public channels in the workspace
func (c *Client) ListChannels(limit int, cursor string) (*GetConversationsResponse, error) {
	params := &slack.GetConversationsParameters{
		Limit:           limit,
		Cursor:          cursor,
		ExcludeArchived: true,
		Types:           []string{"public_channel"},
	}
	channels, nextCursor, err := c.api.GetConversations(params)
	if err != nil {
		return nil, err
	}
	return &GetConversationsResponse{
		Channels: channels,
		ResponseMetadata: struct{ NextCursor string }{
			NextCursor: nextCursor,
		},
	}, nil
}
