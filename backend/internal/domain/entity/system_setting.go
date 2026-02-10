package entity

import "time"

type SystemSetting struct {
	ID             string    `json:"id" gorm:"column:id;primaryKey"`
	SESRegion      string    `json:"ses_region" gorm:"column:ses_region"`
	SESAccessKeyID string    `json:"ses_access_key_id" gorm:"column:ses_access_key_id"`
	SESSecretKey   string    `json:"ses_secret_key" gorm:"column:ses_secret_key"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
