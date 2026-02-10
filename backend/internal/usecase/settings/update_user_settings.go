package settings

import (
	"fmt"
	"strings"

	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type UpdateUserSettingsUseCase struct {
	repo repository.UserSettingRepository
}

func NewUpdateUserSettingsUseCase(repo repository.UserSettingRepository) *UpdateUserSettingsUseCase {
	return &UpdateUserSettingsUseCase{repo: repo}
}

func (uc *UpdateUserSettingsUseCase) Execute(uid string, discordWebhookURL string, selectedDomainID string) error {
	if err := validateDiscordWebhookURL(discordWebhookURL); err != nil {
		return err
	}

	setting := &entity.UserSetting{
		UID:               uid,
		DiscordWebhookURL: strings.TrimSpace(discordWebhookURL),
		SelectedDomainID:  strings.TrimSpace(selectedDomainID),
	}
	return uc.repo.Upsert(setting)
}

func validateDiscordWebhookURL(url string) error {
	if strings.TrimSpace(url) == "" {
		return nil
	}

	validPrefixes := []string{
		"https://discord.com/api/webhooks/",
		"https://discordapp.com/api/webhooks/",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return nil
		}
	}

	return fmt.Errorf("invalid discord webhook url")
}
