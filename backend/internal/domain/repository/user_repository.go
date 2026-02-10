package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type UserRepository interface {
	GetByUID(uid string) (*entity.User, error)
	Upsert(user *entity.User) error
}
