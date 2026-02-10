package entity

import "time"

type UserSetting struct {
	UID               string    `json:"uid" gorm:"column:uid;primaryKey"`
	DiscordWebhookURL string    `json:"discord_webhook_url" gorm:"column:discord_webhook_url"`
	SelectedDomainID  string    `json:"selected_domain_id" gorm:"column:selected_domain_id"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
