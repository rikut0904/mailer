package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rikut0904/mailer-backend/pkg/config"
)

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

type ClientConfig struct {
	FirebaseAPIKey    string `json:"firebase_api_key"`
	FirebaseAuthDomain string `json:"firebase_auth_domain"`
	FirebaseProjectID string `json:"firebase_project_id"`
}

func (h *ConfigHandler) GetClientConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, ClientConfig{
		FirebaseAPIKey:    h.cfg.FirebaseAPIKey,
		FirebaseAuthDomain: h.cfg.FirebaseAuthDomain,
		FirebaseProjectID: h.cfg.FirebaseProjectID,
	})
}
