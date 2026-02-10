package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mailStateRepository struct {
	db *gorm.DB
}

func NewMailStateRepository(db *gorm.DB) repository.MailStateRepository {
	return &mailStateRepository{db: db}
}

func (r *mailStateRepository) FindByS3Key(domainID, s3Key string) (*entity.MailState, error) {
	var state entity.MailState
	if err := r.db.Where("domain_id = ? AND s3_key = ?", domainID, s3Key).First(&state).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *mailStateRepository) FindByRecipient(domainID, recipientAddress string, offset, limit int) ([]entity.MailState, int64, error) {
	var states []entity.MailState
	var total int64

	query := r.db.Model(&entity.MailState{}).Where("domain_id = ?", domainID)
	if recipientAddress != "" {
		query = query.Where("recipient_address = ?", recipientAddress)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&states).Error; err != nil {
		return nil, 0, err
	}

	return states, total, nil
}

func (r *mailStateRepository) FindByThreadID(domainID, threadID string) ([]entity.MailState, error) {
	var states []entity.MailState
	if err := r.db.Where("domain_id = ? AND thread_id = ?", domainID, threadID).Order("created_at ASC").Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (r *mailStateRepository) Upsert(state *entity.MailState) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "domain_id"}, {Name: "s3_key"}},
		DoNothing: true,
	}).Create(state).Error
}

func (r *mailStateRepository) UpdateReadStatus(domainID, s3Key string, isRead bool) error {
	return r.db.Model(&entity.MailState{}).Where("domain_id = ? AND s3_key = ?", domainID, s3Key).Update("is_read", isRead).Error
}

func (r *mailStateRepository) UpdateStarStatus(domainID, s3Key string, isStarred bool) error {
	return r.db.Model(&entity.MailState{}).Where("domain_id = ? AND s3_key = ?", domainID, s3Key).Update("is_starred", isStarred).Error
}

func (r *mailStateRepository) UpdateThreadID(domainID, s3Key string, threadID string) error {
	return r.db.Model(&entity.MailState{}).Where("domain_id = ? AND s3_key = ?", domainID, s3Key).Update("thread_id", threadID).Error
}

func (r *mailStateRepository) Delete(domainID, s3Key string) error {
	return r.db.Where("domain_id = ? AND s3_key = ?", domainID, s3Key).Delete(&entity.MailState{}).Error
}

func (r *mailStateRepository) CountUnread(domainID, recipientAddress string) (int64, error) {
	var count int64
	query := r.db.Model(&entity.MailState{}).Where("domain_id = ? AND is_read = ?", domainID, false)
	if recipientAddress != "" {
		query = query.Where("recipient_address = ?", recipientAddress)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
