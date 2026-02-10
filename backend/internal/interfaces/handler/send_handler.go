package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	senduc "github.com/rikut0904/mailer-backend/internal/usecase/send"
)

type SendHandler struct {
	sendMailUC *senduc.SendMailUseCase
}

func NewSendHandler(sendMailUC *senduc.SendMailUseCase) *SendHandler {
	return &SendHandler{sendMailUC: sendMailUC}
}

func (h *SendHandler) SendMail(c echo.Context) error {
	var req senduc.SendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if len(req.To) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "at least one recipient is required"})
	}
	if req.Subject == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "subject is required"})
	}
	if req.Body == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "body is required"})
	}
	if req.SendType == "" {
		req.SendType = "new"
	}

	result, err := h.sendMailUC.Execute(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
