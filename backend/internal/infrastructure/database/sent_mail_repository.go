package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type sentMailRepository struct {
	db *gorm.DB
}

func NewSentMailRepository(db *gorm.DB) repository.SentMailRepository {
	return &sentMailRepository{db: db}
}

func (r *sentMailRepository) FindByManagementCode(code string) (*entity.SentMail, error) {
	var sentMail entity.SentMail
	if err := r.db.Where("management_code = ?", code).First(&sentMail).Error; err != nil {
		return nil, err
	}
	return &sentMail, nil
}

func (r *sentMailRepository) FindByParentThreadID(threadID string) ([]entity.SentMail, error) {
	var sentMails []entity.SentMail
	if err := r.db.Where("parent_thread_id = ?", threadID).Order("sent_at ASC").Find(&sentMails).Error; err != nil {
		return nil, err
	}
	return sentMails, nil
}

func (r *sentMailRepository) Create(sentMail *entity.SentMail) error {
	return r.db.Create(sentMail).Error
}

func (r *sentMailRepository) Delete(managementCode string) error {
	return r.db.Where("management_code = ?", managementCode).Delete(&entity.SentMail{}).Error
}
