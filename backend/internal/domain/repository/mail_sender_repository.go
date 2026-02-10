package repository

type MailSenderRepository interface {
	SendRawEmail(from, to, subject, textBody, htmlBody string) error
}
