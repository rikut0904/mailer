package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const systemSettingID = "default"

type systemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) repository.SystemSettingRepository {
	return &systemSettingRepository{db: db}
}

func (r *systemSettingRepository) Get() (*entity.SystemSetting, error) {
	var setting entity.SystemSetting
	if err := r.db.Where("id = ?", systemSettingID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *systemSettingRepository) Upsert(setting *entity.SystemSetting) error {
	setting.ID = systemSettingID
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"ses_region", "ses_access_key_id", "ses_secret_key", "updated_at"}),
	}).Create(setting).Error
}
