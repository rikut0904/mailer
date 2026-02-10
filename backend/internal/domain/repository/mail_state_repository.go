package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type MailStateRepository interface {
	FindByS3Key(domainID, s3Key string) (*entity.MailState, error)
	FindByRecipient(domainID, recipientAddress string, offset, limit int) ([]entity.MailState, int64, error)
	FindByThreadID(domainID, threadID string) ([]entity.MailState, error)
	Upsert(state *entity.MailState) error
	UpdateReadStatus(domainID, s3Key string, isRead bool) error
	UpdateStarStatus(domainID, s3Key string, isStarred bool) error
	UpdateThreadID(domainID, s3Key string, threadID string) error
	Delete(domainID, s3Key string) error
	CountUnread(domainID, recipientAddress string) (int64, error)
}
