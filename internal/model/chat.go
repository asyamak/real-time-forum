package model

type Chat struct {
	User                User    `json:"user"`
	LastMessage         Message `json:"last_message"`
	UnreadMessagesCount int     `json:"unread_messages_count"`
}
