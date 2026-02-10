package mail

import (
	"fmt"

	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type UpdateStateUseCase struct {
	mailStateRepo repository.MailStateRepository
}

func NewUpdateStateUseCase(mailStateRepo repository.MailStateRepository) *UpdateStateUseCase {
	return &UpdateStateUseCase{mailStateRepo: mailStateRepo}
}

func (uc *UpdateStateUseCase) MarkAsRead(domainID, s3Key string, isRead bool) error {
	if err := uc.mailStateRepo.UpdateReadStatus(domainID, s3Key, isRead); err != nil {
		return fmt.Errorf("failed to update read status: %w", err)
	}
	return nil
}

func (uc *UpdateStateUseCase) MarkAsStarred(domainID, s3Key string, isStarred bool) error {
	if err := uc.mailStateRepo.UpdateStarStatus(domainID, s3Key, isStarred); err != nil {
		return fmt.Errorf("failed to update star status: %w", err)
	}
	return nil
}
