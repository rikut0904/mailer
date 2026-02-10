package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type UserSettingRepository interface {
	GetByUID(uid string) (*entity.UserSetting, error)
	Upsert(setting *entity.UserSetting) error
}
