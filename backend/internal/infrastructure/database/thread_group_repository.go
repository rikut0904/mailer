package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type threadGroupRepository struct {
	db *gorm.DB
}

func NewThreadGroupRepository(db *gorm.DB) repository.ThreadGroupRepository {
	return &threadGroupRepository{db: db}
}

func (r *threadGroupRepository) FindByParentUUID(parentUUID string) (*entity.ThreadGroup, error) {
	var group entity.ThreadGroup
	if err := r.db.Where("parent_uuid = ?", parentUUID).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *threadGroupRepository) Create(group *entity.ThreadGroup) error {
	return r.db.Create(group).Error
}

func (r *threadGroupRepository) List() ([]entity.ThreadGroup, error) {
	var groups []entity.ThreadGroup
	if err := r.db.Order("parent_uuid").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *threadGroupRepository) Delete(parentUUID string) error {
	return r.db.Where("parent_uuid = ?", parentUUID).Delete(&entity.ThreadGroup{}).Error
}
