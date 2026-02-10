package entity

import "time"

type User struct {
	UID       string    `json:"uid" gorm:"column:uid;primaryKey"`
	Role      string    `json:"role" gorm:"column:role"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}
