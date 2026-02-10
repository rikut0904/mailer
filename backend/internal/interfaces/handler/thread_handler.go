package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	awsinfra "github.com/rikut0904/mailer-backend/internal/infrastructure/aws"
	threaduc "github.com/rikut0904/mailer-backend/internal/usecase/thread"
)

type ThreadHandler struct {
	getThreadUC     *threaduc.GetThreadUseCase
	userSettingRepo repository.UserSettingRepository
	domainRepo      repository.S3DomainRepository
}

func NewThreadHandler(
	getThreadUC *threaduc.GetThreadUseCase,
	userSettingRepo repository.UserSettingRepository,
	domainRepo repository.S3DomainRepository,
) *ThreadHandler {
	return &ThreadHandler{
		getThreadUC:     getThreadUC,
		userSettingRepo: userSettingRepo,
		domainRepo:      domainRepo,
	}
}

func (h *ThreadHandler) GetThread(c echo.Context) error {
	threadID := c.Param("threadId")
	if threadID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "thread_id is required"})
	}

	domain, storageRepo, err := h.storageForUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.getThreadUC.Execute(storageRepo, domain.ID, threadID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ThreadHandler) ListThreads(c echo.Context) error {
	threads, err := h.getThreadUC.ListThreads()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, threads)
}

func (h *ThreadHandler) storageForUser(c echo.Context) (*entity.S3Domain, repository.MailStorageRepository, error) {
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

func (h *ThreadHandler) resolveDomain(uid string, role string) (*entity.S3Domain, error) {
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
