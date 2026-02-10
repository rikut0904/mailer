package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type DomainHandler struct {
	domainRepo repository.S3DomainRepository
}

func NewDomainHandler(domainRepo repository.S3DomainRepository) *DomainHandler {
	return &DomainHandler{domainRepo: domainRepo}
}

type DomainRequest struct {
	Name        string `json:"name"`
	Bucket      string `json:"bucket"`
	Region      string `json:"region"`
	AccessKeyID string `json:"access_key_id"`
	SecretKey   string `json:"secret_key"`
	Endpoint    string `json:"endpoint"`
}

func (h *DomainHandler) ListDomains(c echo.Context) error {
	domains, err := h.domainRepo.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, domains)
}

func (h *DomainHandler) CreateDomain(c echo.Context) error {
	var req DomainRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	domain := &entity.S3Domain{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Bucket:      req.Bucket,
		Region:      req.Region,
		AccessKeyID: req.AccessKeyID,
		SecretKey:   req.SecretKey,
		Endpoint:    req.Endpoint,
	}

	if err := h.domainRepo.Create(domain); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, domain)
}

func (h *DomainHandler) UpdateDomain(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	var req DomainRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	domain, err := h.domainRepo.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "domain not found"})
	}

	domain.Name = req.Name
	domain.Bucket = req.Bucket
	domain.Region = req.Region
	domain.AccessKeyID = req.AccessKeyID
	domain.SecretKey = req.SecretKey
	domain.Endpoint = req.Endpoint

	if err := h.domainRepo.Update(domain); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, domain)
}

func (h *DomainHandler) DeleteDomain(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	if err := h.domainRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
