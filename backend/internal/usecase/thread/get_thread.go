package thread

import (
	"fmt"
	"sort"
	"time"

	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	mimeparser "github.com/rikut0904/mailer-backend/pkg/mime"
)

type GetThreadUseCase struct {
	threadGroupRepo repository.ThreadGroupRepository
	mailStateRepo   repository.MailStateRepository
	sentMailRepo    repository.SentMailRepository
}

func NewGetThreadUseCase(
	threadGroupRepo repository.ThreadGroupRepository,
	mailStateRepo repository.MailStateRepository,
	sentMailRepo repository.SentMailRepository,
) *GetThreadUseCase {
	return &GetThreadUseCase{
		threadGroupRepo: threadGroupRepo,
		mailStateRepo:   mailStateRepo,
		sentMailRepo:    sentMailRepo,
	}
}

type ThreadMessage struct {
	Type      string    `json:"type"` // "received" or "sent"
	Subject   string    `json:"subject"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Body      string    `json:"body"`
	Date      time.Time `json:"date"`
	S3Key     string    `json:"s3_key,omitempty"`
	Code      string    `json:"management_code,omitempty"`
	IsRead    bool      `json:"is_read,omitempty"`
	IsStarred bool      `json:"is_starred,omitempty"`
}

type ThreadResponse struct {
	ThreadID  string          `json:"thread_id"`
	GroupName string          `json:"group_name"`
	Messages  []ThreadMessage `json:"messages"`
}

func (uc *GetThreadUseCase) Execute(storageRepo repository.MailStorageRepository, domainID, threadID string) (*ThreadResponse, error) {
	group, err := uc.threadGroupRepo.FindByParentUUID(threadID)
	if err != nil {
		return nil, fmt.Errorf("thread group not found: %w", err)
	}

	var messages []ThreadMessage

	receivedMails, err := uc.mailStateRepo.FindByThreadID(domainID, threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch received mails: %w", err)
	}

	for _, state := range receivedMails {
		raw, err := storageRepo.GetObject(state.S3Key)
		if err != nil {
			continue
		}
		parsed, err := mimeparser.Parse(raw, state.S3Key)
		if err != nil {
			continue
		}
		messages = append(messages, ThreadMessage{
			Type:      "received",
			Subject:   parsed.Subject,
			From:      parsed.From,
			To:        parsed.To,
			Body:      parsed.Body,
			Date:      parsed.Date,
			S3Key:     state.S3Key,
			IsRead:    state.IsRead,
			IsStarred: state.IsStarred,
		})
	}

	sentMails, err := uc.sentMailRepo.FindByParentThreadID(threadID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sent mails: %w", err)
	}

	for _, sent := range sentMails {
		messages = append(messages, ThreadMessage{
			Type:    "sent",
			Subject: sent.Subject,
			To:      sent.RecipientEmail,
			Body:    sent.Body,
			Date:    sent.SentAt,
			Code:    sent.ManagementCode,
		})
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Date.Before(messages[j].Date)
	})

	return &ThreadResponse{
		ThreadID:  group.ParentUUID,
		GroupName: group.GroupName,
		Messages:  messages,
	}, nil
}

func (uc *GetThreadUseCase) ListThreads() ([]entity.ThreadGroup, error) {
	return uc.threadGroupRepo.List()
}
