package models

import (
	"bytes"
)

type contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type MsgRequest struct {
	ChatId    string        `json:"chat_id"`
	Text      string        `json:"text,omitempty"`
	Image     *bytes.Buffer `json:"image,omitempty"`
	Audio     *bytes.Buffer `json:"audio,omitempty"`
	Document  *bytes.Buffer `json:"document,omitempty"`
	Video     *bytes.Buffer `json:"video,omitempty"`
	Animation *bytes.Buffer `json:"animation,omitempty"`
	Voice     *bytes.Buffer `json:"voice,omitempty"`
	Contact   *contact      `json:"contact,omitempty"`
	Caption   string        `json:"caption,omitempty"`
}
