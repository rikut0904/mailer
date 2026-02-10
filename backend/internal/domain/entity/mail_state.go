package entity

import "time"

type MailState struct {
	S3Key            string    `json:"s3_key" gorm:"column:s3_key;primaryKey"`
	DomainID         string    `json:"domain_id" gorm:"column:domain_id;primaryKey"`
	RecipientAddress string    `json:"recipient_address" gorm:"column:recipient_address"`
	IsRead           bool      `json:"is_read" gorm:"column:is_read;default:false"`
	IsStarred        bool      `json:"is_starred" gorm:"column:is_starred;default:false"`
	ThreadID         *string   `json:"thread_id,omitempty" gorm:"column:thread_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (MailState) TableName() string {
	return "mail_states"
}
