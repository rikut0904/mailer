package mail

import (
	"fmt"
	"log"
	"strings"

	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	mimeparser "github.com/rikut0904/mailer-backend/pkg/mime"
)

type SyncMailsUseCase struct {
	mailStateRepo repository.MailStateRepository
	threadLinkUC  *LinkThreadUseCase
}

func NewSyncMailsUseCase(
	mailStateRepo repository.MailStateRepository,
	threadLinkUC *LinkThreadUseCase,
) *SyncMailsUseCase {
	return &SyncMailsUseCase{
		mailStateRepo: mailStateRepo,
		threadLinkUC:  threadLinkUC,
	}
}

func (uc *SyncMailsUseCase) Execute(storageRepo repository.MailStorageRepository, domainID string) (int, error) {
	var synced int
	var continuationToken *string

	for {
		keys, nextToken, err := storageRepo.ListKeys("", continuationToken, 100)
		if err != nil {
			return synced, fmt.Errorf("failed to list S3 keys: %w", err)
		}

		for _, key := range keys {
			if strings.HasSuffix(key, "/") {
				continue
			}

			existing, _ := uc.mailStateRepo.FindByS3Key(domainID, key)
			if existing != nil {
				continue
			}

			raw, err := storageRepo.GetObject(key)
			if err != nil {
				log.Printf("failed to get S3 object %s: %v", key, err)
				continue
			}

			parsed, err := mimeparser.Parse(raw, key)
			if err != nil {
				log.Printf("failed to parse mail %s: %v", key, err)
				continue
			}

			state := &entity.MailState{
				S3Key:            key,
				DomainID:         domainID,
				RecipientAddress: parsed.To,
				IsRead:           false,
				IsStarred:        false,
				CreatedAt:        parsed.Date,
			}

			if err := uc.mailStateRepo.Upsert(state); err != nil {
				log.Printf("failed to upsert mail state %s: %v", key, err)
				continue
			}

			if uc.threadLinkUC != nil {
				if threadID, err := uc.threadLinkUC.LinkFromBody(parsed.Body, domainID, key); err == nil && threadID != "" {
					log.Printf("linked mail %s to thread %s", key, threadID)
				}
			}

			synced++
		}

		if nextToken == nil {
			break
		}
		continuationToken = nextToken
	}

	return synced, nil
}
