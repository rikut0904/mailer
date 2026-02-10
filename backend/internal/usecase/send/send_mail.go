package send

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"github.com/rikut0904/mailer-backend/internal/infrastructure/discord"
)

type SendMailUseCase struct {
	sentMailRepo    repository.SentMailRepository
	threadGroupRepo repository.ThreadGroupRepository
	senderRepo      repository.MailSenderRepository
	discordClient   *discord.Client
}

func NewSendMailUseCase(
	sentMailRepo repository.SentMailRepository,
	threadGroupRepo repository.ThreadGroupRepository,
	senderRepo repository.MailSenderRepository,
	discordClient *discord.Client,
) *SendMailUseCase {
	return &SendMailUseCase{
		sentMailRepo:    sentMailRepo,
		threadGroupRepo: threadGroupRepo,
		senderRepo:      senderRepo,
		discordClient:   discordClient,
	}
}

type SendRequest struct {
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	HTMLBody    string   `json:"html_body,omitempty"`
	ThreadID    string   `json:"thread_id,omitempty"`
	ReplyCode   string   `json:"reply_code,omitempty"`
	SendType    string   `json:"send_type"` // "new", "reply", "forward"
	FromAddress string   `json:"from_address,omitempty"`
}

type SendResponse struct {
	ThreadID        string   `json:"thread_id"`
	ManagementCodes []string `json:"management_codes"`
}

func (uc *SendMailUseCase) Execute(req *SendRequest) (*SendResponse, error) {
	var threadID string
	var managementCodes []string

	from := req.FromAddress
	if strings.TrimSpace(from) == "" {
		return nil, fmt.Errorf("from_address is required")
	}

	switch req.SendType {
	case "new":
		threadID = uuid.New().String()
		if err := uc.threadGroupRepo.Create(&entity.ThreadGroup{
			ParentUUID: threadID,
			GroupName:  req.Subject,
		}); err != nil {
			return nil, fmt.Errorf("failed to create thread group: %w", err)
		}

		for _, to := range req.To {
			code := uuid.New().String()
			managementCodes = append(managementCodes, code)

			bodyWithCode := appendManagementCode(req.Body, code)
			htmlBodyWithCode := ""
			if req.HTMLBody != "" {
				htmlBodyWithCode = appendManagementCodeHTML(req.HTMLBody, code)
			}

			if err := uc.senderRepo.SendRawEmail(from, to, req.Subject, bodyWithCode, htmlBodyWithCode); err != nil {
				return nil, fmt.Errorf("failed to send email to %s: %w", to, err)
			}

			if err := uc.sentMailRepo.Create(&entity.SentMail{
				ManagementCode: code,
				ParentThreadID: threadID,
				RecipientEmail: to,
				Subject:        req.Subject,
				Body:           bodyWithCode,
			}); err != nil {
				return nil, fmt.Errorf("failed to save sent mail record: %w", err)
			}
		}

	case "reply":
		if req.ThreadID == "" {
			return nil, fmt.Errorf("thread_id is required for reply")
		}
		threadID = req.ThreadID

		var code string
		if req.ReplyCode != "" {
			code = req.ReplyCode
		} else {
			code = uuid.New().String()
		}
		managementCodes = append(managementCodes, code)

		for _, to := range req.To {
			bodyWithCode := appendManagementCode(req.Body, code)
			htmlBodyWithCode := ""
			if req.HTMLBody != "" {
				htmlBodyWithCode = appendManagementCodeHTML(req.HTMLBody, code)
			}

			if err := uc.senderRepo.SendRawEmail(from, to, req.Subject, bodyWithCode, htmlBodyWithCode); err != nil {
				return nil, fmt.Errorf("failed to send reply to %s: %w", to, err)
			}

			if err := uc.sentMailRepo.Create(&entity.SentMail{
				ManagementCode: code,
				ParentThreadID: threadID,
				RecipientEmail: to,
				Subject:        req.Subject,
				Body:           bodyWithCode,
			}); err != nil {
				return nil, fmt.Errorf("failed to save sent mail record: %w", err)
			}
		}

	case "forward":
		if req.ThreadID == "" {
			return nil, fmt.Errorf("thread_id is required for forward")
		}
		threadID = req.ThreadID

		if len(req.To) == 1 {
			code := uuid.New().String()
			managementCodes = append(managementCodes, code)

			bodyWithCode := appendManagementCode(req.Body, code)
			htmlBodyWithCode := ""
			if req.HTMLBody != "" {
				htmlBodyWithCode = appendManagementCodeHTML(req.HTMLBody, code)
			}

			if err := uc.senderRepo.SendRawEmail(from, req.To[0], req.Subject, bodyWithCode, htmlBodyWithCode); err != nil {
				return nil, fmt.Errorf("failed to forward email: %w", err)
			}

			if err := uc.sentMailRepo.Create(&entity.SentMail{
				ManagementCode: code,
				ParentThreadID: threadID,
				RecipientEmail: req.To[0],
				Subject:        req.Subject,
				Body:           bodyWithCode,
			}); err != nil {
				return nil, fmt.Errorf("failed to save sent mail record: %w", err)
			}
		} else {
			for _, to := range req.To {
				childCode := uuid.New().String()
				managementCodes = append(managementCodes, childCode)

				bodyWithCode := appendManagementCode(req.Body, childCode)
				htmlBodyWithCode := ""
				if req.HTMLBody != "" {
					htmlBodyWithCode = appendManagementCodeHTML(req.HTMLBody, childCode)
				}

				if err := uc.senderRepo.SendRawEmail(from, to, req.Subject, bodyWithCode, htmlBodyWithCode); err != nil {
					return nil, fmt.Errorf("failed to forward email to %s: %w", to, err)
				}

				if err := uc.sentMailRepo.Create(&entity.SentMail{
					ManagementCode: childCode,
					ParentThreadID: threadID,
					RecipientEmail: to,
					Subject:        req.Subject,
					Body:           bodyWithCode,
				}); err != nil {
					return nil, fmt.Errorf("failed to save sent mail record: %w", err)
				}
			}
		}

	default:
		return nil, fmt.Errorf("unsupported send type: %s", req.SendType)
	}

	return &SendResponse{
		ThreadID:        threadID,
		ManagementCodes: managementCodes,
	}, nil
}

func appendManagementCode(body, code string) string {
	signature := fmt.Sprintf("\n\n---\n【管理コード: %s】", code)
	return strings.TrimRight(body, "\n") + signature
}

func appendManagementCodeHTML(html, code string) string {
	signature := fmt.Sprintf(`<br><hr><p style="font-size:12px;color:#888;">【管理コード: %s】</p>`, code)
	if idx := strings.LastIndex(html, "</body>"); idx != -1 {
		return html[:idx] + signature + html[idx:]
	}
	return html + signature
}
