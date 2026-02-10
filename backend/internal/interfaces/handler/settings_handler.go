package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	settingsuc "github.com/rikut0904/mailer-backend/internal/usecase/settings"
)

type SettingsHandler struct {
	getUC    *settingsuc.GetUserSettingsUseCase
	updateUC *settingsuc.UpdateUserSettingsUseCase
}

func NewSettingsHandler(
	getUC *settingsuc.GetUserSettingsUseCase,
	updateUC *settingsuc.UpdateUserSettingsUseCase,
) *SettingsHandler {
	return &SettingsHandler{
		getUC:    getUC,
		updateUC: updateUC,
	}
}

type UpdateSettingsRequest struct {
	DiscordWebhookURL string `json:"discord_webhook_url"`
	SelectedDomainID  string `json:"selected_domain_id"`
}

type SettingsResponse struct {
	DiscordWebhookURL string `json:"discord_webhook_url"`
	SelectedDomainID  string `json:"selected_domain_id"`
}

func (h *SettingsHandler) GetSettings(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	setting, err := h.getUC.Execute(uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, SettingsResponse{
		DiscordWebhookURL: setting.DiscordWebhookURL,
		SelectedDomainID:  setting.SelectedDomainID,
	})
}

func (h *SettingsHandler) UpdateSettings(c echo.Context) error {
	uid, ok := c.Get("uid").(string)
	if !ok || uid == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req UpdateSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.updateUC.Execute(uid, req.DiscordWebhookURL, req.SelectedDomainID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, SettingsResponse{
		DiscordWebhookURL: req.DiscordWebhookURL,
		SelectedDomainID:  req.SelectedDomainID,
	})
}
