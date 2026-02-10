package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type SystemSettingHandler struct {
	repo repository.SystemSettingRepository
}

func NewSystemSettingHandler(repo repository.SystemSettingRepository) *SystemSettingHandler {
	return &SystemSettingHandler{repo: repo}
}

type SystemSettingRequest struct {
	SESRegion      string `json:"ses_region"`
	SESAccessKeyID string `json:"ses_access_key_id"`
	SESSecretKey   string `json:"ses_secret_key"`
}

func (h *SystemSettingHandler) Get(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	setting, err := h.repo.Get()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "system settings not found"})
	}

	return c.JSON(http.StatusOK, setting)
}

func (h *SystemSettingHandler) Update(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	var req SystemSettingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	setting := &entity.SystemSetting{
		SESRegion:      req.SESRegion,
		SESAccessKeyID: req.SESAccessKeyID,
		SESSecretKey:   req.SESSecretKey,
	}

	if err := h.repo.Upsert(setting); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, setting)
}
