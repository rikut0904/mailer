package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type SystemSettingRepository interface {
	Get() (*entity.SystemSetting, error)
	Upsert(setting *entity.SystemSetting) error
}
