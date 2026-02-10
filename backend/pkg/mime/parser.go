package mime

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"time"

	"github.com/rikut0904/mailer-backend/internal/domain/entity"
)

func Parse(raw []byte, s3Key string) (*entity.ParsedMail, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	parsed := &entity.ParsedMail{
		S3Key:     s3Key,
		MessageID: msg.Header.Get("Message-ID"),
		From:      decodeHeader(msg.Header.Get("From")),
		To:        decodeHeader(msg.Header.Get("To")),
		Subject:   decodeHeader(msg.Header.Get("Subject")),
	}

	if dateStr := msg.Header.Get("Date"); dateStr != "" {
		if t, err := mail.ParseDate(dateStr); err == nil {
			parsed.Date = t
		}
	}
	if parsed.Date.IsZero() {
		parsed.Date = time.Now()
	}

	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain"
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		body, _ := io.ReadAll(msg.Body)
		parsed.Body = string(body)
		return parsed, nil
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		err = parseMultipart(msg.Body, params["boundary"], parsed)
		if err != nil {
			return nil, fmt.Errorf("failed to parse multipart: %w", err)
		}
	} else {
		body, err := readBody(msg.Body, msg.Header.Get("Content-Transfer-Encoding"))
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		if strings.Contains(mediaType, "html") {
			parsed.HTMLBody = body
		} else {
			parsed.Body = body
		}
	}

	return parsed, nil
}

func parseMultipart(r io.Reader, boundary string, parsed *entity.ParsedMail) error {
	mr := multipart.NewReader(r, boundary)
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		partContentType := part.Header.Get("Content-Type")
		mediaType, params, _ := mime.ParseMediaType(partContentType)
		disposition := part.Header.Get("Content-Disposition")

		if strings.HasPrefix(mediaType, "multipart/") {
			if err := parseMultipart(part, params["boundary"], parsed); err != nil {
				return err
			}
			continue
		}

		if strings.Contains(disposition, "attachment") || part.FileName() != "" {
			data, err := readBodyBytes(part, part.Header.Get("Content-Transfer-Encoding"))
			if err != nil {
				continue
			}
			parsed.Attachments = append(parsed.Attachments, entity.Attachment{
				Filename:    part.FileName(),
				ContentType: mediaType,
				Size:        len(data),
				Content:     data,
			})
			continue
		}

		body, err := readBody(part, part.Header.Get("Content-Transfer-Encoding"))
		if err != nil {
			continue
		}

		if strings.Contains(mediaType, "html") {
			if parsed.HTMLBody == "" {
				parsed.HTMLBody = body
			}
		} else if strings.Contains(mediaType, "plain") {
			if parsed.Body == "" {
				parsed.Body = body
			}
		}
	}
	return nil
}

func readBody(r io.Reader, encoding string) (string, error) {
	data, err := readBodyBytes(r, encoding)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readBodyBytes(r io.Reader, encoding string) ([]byte, error) {
	switch strings.ToLower(encoding) {
	case "base64":
		decoder := base64.NewDecoder(base64.StdEncoding, r)
		return io.ReadAll(decoder)
	case "quoted-printable":
		decoder := quotedprintable.NewReader(r)
		return io.ReadAll(decoder)
	default:
		return io.ReadAll(r)
	}
}

func decodeHeader(s string) string {
	dec := new(mime.WordDecoder)
	decoded, err := dec.DecodeHeader(s)
	if err != nil {
		return s
	}
	return decoded
}
