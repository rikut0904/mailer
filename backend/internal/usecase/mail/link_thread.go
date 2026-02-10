package mail

import (
	"regexp"

	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

var managementCodeRegex = regexp.MustCompile(`【管理コード: ([^】]+)】`)

type LinkThreadUseCase struct {
	sentMailRepo  repository.SentMailRepository
	mailStateRepo repository.MailStateRepository
}

func NewLinkThreadUseCase(
	sentMailRepo repository.SentMailRepository,
	mailStateRepo repository.MailStateRepository,
) *LinkThreadUseCase {
	return &LinkThreadUseCase{
		sentMailRepo:  sentMailRepo,
		mailStateRepo: mailStateRepo,
	}
}

func (uc *LinkThreadUseCase) LinkFromBody(body string, domainID string, s3Key string) (string, error) {
	matches := managementCodeRegex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return "", nil
	}

	code := matches[1]
	sentMail, err := uc.sentMailRepo.FindByManagementCode(code)
	if err != nil {
		return "", nil
	}

	if err := uc.mailStateRepo.UpdateThreadID(domainID, s3Key, sentMail.ParentThreadID); err != nil {
		return "", err
	}

	return sentMail.ParentThreadID, nil
}

func ExtractManagementCode(body string) string {
	matches := managementCodeRegex.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
