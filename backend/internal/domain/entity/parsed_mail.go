package entity

import "time"

type ParsedMail struct {
	S3Key       string       `json:"s3_key"`
	MessageID   string       `json:"message_id"`
	From        string       `json:"from"`
	To          string       `json:"to"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	HTMLBody    string       `json:"html_body,omitempty"`
	Date        time.Time    `json:"date"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// State from DB
	IsRead    bool    `json:"is_read"`
	IsStarred bool    `json:"is_starred"`
	ThreadID  *string `json:"thread_id,omitempty"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	Content     []byte `json:"-"`
}
