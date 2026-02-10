package settings

import (
	"errors"

	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type GetUserSettingsUseCase struct {
	repo repository.UserSettingRepository
}

func NewGetUserSettingsUseCase(repo repository.UserSettingRepository) *GetUserSettingsUseCase {
	return &GetUserSettingsUseCase{repo: repo}
}

func (uc *GetUserSettingsUseCase) Execute(uid string) (*entity.UserSetting, error) {
	setting, err := uc.repo.GetByUID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entity.UserSetting{
				UID:               uid,
				DiscordWebhookURL: "",
			}, nil
		}
		return nil, err
	}
	return setting, nil
}
