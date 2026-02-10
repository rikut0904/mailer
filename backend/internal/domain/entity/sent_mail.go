package entity

import "time"

type SentMail struct {
	ManagementCode string    `json:"management_code" gorm:"column:management_code;primaryKey"`
	ParentThreadID string    `json:"parent_thread_id" gorm:"column:parent_thread_id"`
	RecipientEmail string    `json:"recipient_email" gorm:"column:recipient_email"`
	Subject        string    `json:"subject" gorm:"column:subject"`
	Body           string    `json:"body" gorm:"column:body;type:text"`
	SentAt         time.Time `json:"sent_at" gorm:"column:sent_at;autoCreateTime"`
}

func (SentMail) TableName() string {
	return "sent_mails"
}
