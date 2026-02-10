package database

import (
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"gorm.io/gorm"
)

type s3DomainRepository struct {
	db *gorm.DB
}

func NewS3DomainRepository(db *gorm.DB) repository.S3DomainRepository {
	return &s3DomainRepository{db: db}
}

func (r *s3DomainRepository) List() ([]entity.S3Domain, error) {
	var domains []entity.S3Domain
	if err := r.db.Order("created_at asc").Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

func (r *s3DomainRepository) GetByID(id string) (*entity.S3Domain, error) {
	var domain entity.S3Domain
	if err := r.db.Where("id = ?", id).First(&domain).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

func (r *s3DomainRepository) Create(domain *entity.S3Domain) error {
	return r.db.Create(domain).Error
}

func (r *s3DomainRepository) Update(domain *entity.S3Domain) error {
	return r.db.Save(domain).Error
}

func (r *s3DomainRepository) Delete(id string) error {
	return r.db.Delete(&entity.S3Domain{}, "id = ?", id).Error
}
