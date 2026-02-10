package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.MailState{},
		&entity.ThreadGroup{},
		&entity.SentMail{},
		&entity.UserSetting{},
		&entity.User{},
		&entity.S3Domain{},
		&entity.SystemSetting{},
	)
}
