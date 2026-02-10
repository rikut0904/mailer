package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type SentMailRepository interface {
	FindByManagementCode(code string) (*entity.SentMail, error)
	FindByParentThreadID(threadID string) ([]entity.SentMail, error)
	Create(sentMail *entity.SentMail) error
	Delete(managementCode string) error
}
