package models

type ContactRequest struct {
	ChatId      string `json:"chat_id"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type MsgRequest struct {
	ChatId    string   `json:"chat_id"`
	Text      string   `json:"text,omitempty"`
	ReplyID   string   `json:"reply_id,omitempty"`
	MentionID []string `json:"mention_id,omitempty"`
}
