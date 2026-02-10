package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	awsinfra "github.com/rikut0904/mailer-backend/internal/infrastructure/aws"
	mailuc "github.com/rikut0904/mailer-backend/internal/usecase/mail"
	mimeparser "github.com/rikut0904/mailer-backend/pkg/mime"
)

type MailHandler struct {
	getMailsUC      *mailuc.GetMailsUseCase
	updateStateUC   *mailuc.UpdateStateUseCase
	deleteMailUC    *mailuc.DeleteMailUseCase
	syncMailsUC     *mailuc.SyncMailsUseCase
	userSettingRepo repository.UserSettingRepository
	domainRepo      repository.S3DomainRepository
}

func NewMailHandler(
	getMailsUC *mailuc.GetMailsUseCase,
	updateStateUC *mailuc.UpdateStateUseCase,
	deleteMailUC *mailuc.DeleteMailUseCase,
	syncMailsUC *mailuc.SyncMailsUseCase,
	userSettingRepo repository.UserSettingRepository,
	domainRepo repository.S3DomainRepository,
) *MailHandler {
	return &MailHandler{
		getMailsUC:      getMailsUC,
		updateStateUC:   updateStateUC,
		deleteMailUC:    deleteMailUC,
		syncMailsUC:     syncMailsUC,
		userSettingRepo: userSettingRepo,
		domainRepo:      domainRepo,
	}
}

func (h *MailHandler) GetMails(c echo.Context) error {
	recipient := c.QueryParam("recipient")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	domain, storageRepo, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.getMailsUC.Execute(storageRepo, domain.ID, recipient, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *MailHandler) GetMail(c echo.Context) error {
	s3Key := c.Param("s3Key")
	if s3Key == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "s3_key is required"})
	}

	domain, storageRepo, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	mail, err := h.getMailsUC.GetByS3Key(storageRepo, domain.ID, s3Key)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mail)
}

type UpdateReadRequest struct {
	IsRead bool `json:"is_read"`
}

func (h *MailHandler) UpdateReadStatus(c echo.Context) error {
	s3Key := c.Param("s3Key")
	var req UpdateReadRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	domain, _, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.updateStateUC.MarkAsRead(domain.ID, s3Key, req.IsRead); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

type UpdateStarRequest struct {
	IsStarred bool `json:"is_starred"`
}

func (h *MailHandler) UpdateStarStatus(c echo.Context) error {
	s3Key := c.Param("s3Key")
	var req UpdateStarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	domain, _, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.updateStateUC.MarkAsStarred(domain.ID, s3Key, req.IsStarred); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *MailHandler) DeleteMail(c echo.Context) error {
	s3Key := c.Param("s3Key")

	domain, storageRepo, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.deleteMailUC.Execute(storageRepo, domain.ID, s3Key); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *MailHandler) SyncMails(c echo.Context) error {
	domain, storageRepo, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	count, err := h.syncMailsUC.Execute(storageRepo, domain.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"synced": count,
	})
}

func (h *MailHandler) GetRecipients(c echo.Context) error {
	_, storageRepo, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	keys, _, err := storageRepo.ListKeys("", nil, 200)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	unique := map[string]struct{}{}
	for _, key := range keys {
		if strings.HasSuffix(key, "/") {
			continue
		}

		raw, err := storageRepo.GetObject(key)
		if err != nil {
			continue
		}

		parsed, err := mimeparser.Parse(raw, key)
		if err != nil {
			continue
		}

		for _, addr := range splitRecipients(parsed.To) {
			if addr != "" {
				unique[addr] = struct{}{}
			}
		}
	}

	recipients := make([]string, 0, len(unique))
	for addr := range unique {
		recipients = append(recipients, addr)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recipients": recipients,
	})
}

func splitRecipients(to string) []string {
	parts := strings.Split(to, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (h *MailHandler) storageForUser(c echo.Context) (*entity.S3Domain, repository.MailStorageRepository, error) {
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return nil, nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	role, _ := c.Get("role").(string)

	domain, err := h.resolveDomain(uid, role)
	if err != nil {
		return nil, nil, err
	}

	storageRepo, err := awsinfra.NewS3ClientFromDomain(domain)
	if err != nil {
		return nil, nil, err
	}
	return domain, storageRepo, nil
}

func (h *MailHandler) resolveDomain(uid string, role string) (*entity.S3Domain, error) {
	setting, err := h.userSettingRepo.GetByUID(uid)
	if err == nil && setting.SelectedDomainID != "" {
		if domain, err := h.domainRepo.GetByID(setting.SelectedDomainID); err == nil {
			return domain, nil
		}
	}

	domains, err := h.domainRepo.List()
	if err != nil {
		return nil, err
	}
	if len(domains) == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "no s3 domain configured")
	}

	return &domains[0], nil
}
