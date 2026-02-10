package repository

import "github.com/rikut0904/mailer-backend/internal/domain/entity"

type S3DomainRepository interface {
	List() ([]entity.S3Domain, error)
	GetByID(id string) (*entity.S3Domain, error)
	Create(domain *entity.S3Domain) error
	Update(domain *entity.S3Domain) error
	Delete(id string) error
}
