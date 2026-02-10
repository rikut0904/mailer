package mail

import (
	"fmt"

	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type DeleteMailUseCase struct {
	mailStateRepo repository.MailStateRepository
}

func NewDeleteMailUseCase(
	mailStateRepo repository.MailStateRepository,
) *DeleteMailUseCase {
	return &DeleteMailUseCase{
		mailStateRepo: mailStateRepo,
	}
}

func (uc *DeleteMailUseCase) Execute(storageRepo repository.MailStorageRepository, domainID, s3Key string) error {
	if err := storageRepo.DeleteObject(s3Key); err != nil {
		return fmt.Errorf("failed to delete S3 object: %w", err)
	}

	if err := uc.mailStateRepo.Delete(domainID, s3Key); err != nil {
		return fmt.Errorf("failed to delete mail state: %w", err)
	}

	return nil
}
