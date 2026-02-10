package mail

import (
	"fmt"
	"log"

	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	mimeparser "github.com/rikut0904/mailer-backend/pkg/mime"
)

type GetMailsUseCase struct {
	mailStateRepo repository.MailStateRepository
	threadLinkUC  *LinkThreadUseCase
}

func NewGetMailsUseCase(
	mailStateRepo repository.MailStateRepository,
	threadLinkUC *LinkThreadUseCase,
) *GetMailsUseCase {
	return &GetMailsUseCase{
		mailStateRepo: mailStateRepo,
		threadLinkUC:  threadLinkUC,
	}
}

type MailListResponse struct {
	Mails      []entity.ParsedMail `json:"mails"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalPages int                 `json:"total_pages"`
}

func (uc *GetMailsUseCase) Execute(storageRepo repository.MailStorageRepository, domainID, recipientAddress string, page, perPage int) (*MailListResponse, error) {
	if perPage <= 0 {
		perPage = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage

	states, total, err := uc.mailStateRepo.FindByRecipient(domainID, recipientAddress, offset, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mail states: %w", err)
	}

	mails := make([]entity.ParsedMail, 0, len(states))
	for _, state := range states {
		raw, err := storageRepo.GetObject(state.S3Key)
		if err != nil {
			log.Printf("failed to get S3 object %s: %v", state.S3Key, err)
			continue
		}

		parsed, err := mimeparser.Parse(raw, state.S3Key)
		if err != nil {
			log.Printf("failed to parse mail %s: %v", state.S3Key, err)
			continue
		}

		parsed.IsRead = state.IsRead
		parsed.IsStarred = state.IsStarred
		parsed.ThreadID = state.ThreadID
		mails = append(mails, *parsed)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &MailListResponse{
		Mails:      mails,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (uc *GetMailsUseCase) GetByS3Key(storageRepo repository.MailStorageRepository, domainID, s3Key string) (*entity.ParsedMail, error) {
	state, err := uc.mailStateRepo.FindByS3Key(domainID, s3Key)
	if err != nil {
		return nil, fmt.Errorf("mail state not found: %w", err)
	}

	raw, err := storageRepo.GetObject(s3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get S3 object: %w", err)
	}

	parsed, err := mimeparser.Parse(raw, s3Key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mail: %w", err)
	}

	parsed.IsRead = state.IsRead
	parsed.IsStarred = state.IsStarred
	parsed.ThreadID = state.ThreadID

	return parsed, nil
}
