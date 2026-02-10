package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) repository.UserSettingRepository {
	return &userSettingRepository{db: db}
}

func (r *userSettingRepository) GetByUID(uid string) (*entity.UserSetting, error) {
	var setting entity.UserSetting
	if err := r.db.Where("uid = ?", uid).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *userSettingRepository) Upsert(setting *entity.UserSetting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		DoUpdates: clause.AssignmentColumns([]string{"discord_webhook_url", "selected_domain_id", "updated_at"}),
	}).Create(setting).Error
}
