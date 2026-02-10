package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type ThreadGroupRepository interface {
	FindByParentUUID(parentUUID string) (*entity.ThreadGroup, error)
	Create(group *entity.ThreadGroup) error
	List() ([]entity.ThreadGroup, error)
	Delete(parentUUID string) error
}
