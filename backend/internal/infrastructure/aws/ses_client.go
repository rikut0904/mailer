package aws

import (
	"context"
	"fmt"
	"mime"
	"net/mail"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type sesClient struct {
	settingsRepo repository.SystemSettingRepository
}

func NewSESClient(settingsRepo repository.SystemSettingRepository) repository.MailSenderRepository {
	return &sesClient{settingsRepo: settingsRepo}
}

func (s *sesClient) SendRawEmail(from, to, subject, textBody, htmlBody string) error {
	if strings.TrimSpace(from) == "" {
		return fmt.Errorf("from is required")
	}

	settings, err := s.settingsRepo.Get()
	if err != nil {
		return fmt.Errorf("failed to load SES settings: %w", err)
	}
	if settings.SESRegion == "" || settings.SESAccessKeyID == "" || settings.SESSecretKey == "" {
		return fmt.Errorf("SES settings are not configured")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(settings.SESRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			settings.SESAccessKeyID,
			settings.SESSecretKey,
			"",
		)),
	)
	if err != nil {
		return err
	}
	client := sesv2.NewFromConfig(awsCfg)

	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)

	boundary := fmt.Sprintf("boundary_%d", time.Now().UnixNano())
	var rawMsg strings.Builder

	rawMsg.WriteString(fmt.Sprintf("From: %s\r\n", (&mail.Address{Address: from}).String()))
	rawMsg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	rawMsg.WriteString(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	rawMsg.WriteString("MIME-Version: 1.0\r\n")
	rawMsg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	rawMsg.WriteString("\r\n")

	rawMsg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	rawMsg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	rawMsg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	rawMsg.WriteString(textBody)
	rawMsg.WriteString("\r\n")

	if htmlBody != "" {
		rawMsg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		rawMsg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		rawMsg.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		rawMsg.WriteString(htmlBody)
		rawMsg.WriteString("\r\n")
	}

	rawMsg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	_, err = client.SendEmail(context.TODO(), &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: []byte(rawMsg.String()),
			},
		},
	})

	return err
}
