package entity

import "time"

type S3Domain struct {
	ID          string    `json:"id" gorm:"column:id;primaryKey"`
	Name        string    `json:"name" gorm:"column:name"`
	Bucket      string    `json:"bucket" gorm:"column:bucket"`
	Region      string    `json:"region" gorm:"column:region"`
	AccessKeyID string    `json:"access_key_id" gorm:"column:access_key_id"`
	SecretKey   string    `json:"secret_key" gorm:"column:secret_key"`
	Endpoint    string    `json:"endpoint" gorm:"column:endpoint"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
