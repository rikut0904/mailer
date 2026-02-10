package router

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	"github.com/rikut0904/mailer-backend/internal/infrastructure/discord"
	fbinfra "github.com/rikut0904/mailer-backend/internal/infrastructure/firebase"
	"github.com/rikut0904/mailer-backend/internal/interfaces/handler"
	"github.com/rikut0904/mailer-backend/internal/interfaces/middleware"
	mailuc "github.com/rikut0904/mailer-backend/internal/usecase/mail"
	senduc "github.com/rikut0904/mailer-backend/internal/usecase/send"
	settingsuc "github.com/rikut0904/mailer-backend/internal/usecase/settings"
	threaduc "github.com/rikut0904/mailer-backend/internal/usecase/thread"
	"github.com/rikut0904/mailer-backend/pkg/config"
)

func NewRouter(
	cfg *config.Config,
	fbAuth *fbinfra.FirebaseAuth,
	userRepo repository.UserRepository,
	mailStateRepo repository.MailStateRepository,
	threadGroupRepo repository.ThreadGroupRepository,
	sentMailRepo repository.SentMailRepository,
	userSettingRepo repository.UserSettingRepository,
	domainRepo repository.S3DomainRepository,
	systemSettingRepo repository.SystemSettingRepository,
	senderRepo repository.MailSenderRepository,
	discordClient *discord.Client,
) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(middleware.CORS(cfg.AllowedOrigins))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Public config (Firebase client config for frontend)
	configHandler := handler.NewConfigHandler(cfg)
	e.GET("/api/config", configHandler.GetClientConfig)

	// Usecases
	linkThreadUC := mailuc.NewLinkThreadUseCase(sentMailRepo, mailStateRepo)
	getMailsUC := mailuc.NewGetMailsUseCase(mailStateRepo, linkThreadUC)
	updateStateUC := mailuc.NewUpdateStateUseCase(mailStateRepo)
	deleteMailUC := mailuc.NewDeleteMailUseCase(mailStateRepo)
	syncMailsUC := mailuc.NewSyncMailsUseCase(mailStateRepo, linkThreadUC)
	getThreadUC := threaduc.NewGetThreadUseCase(threadGroupRepo, mailStateRepo, sentMailRepo)
	sendMailUC := senduc.NewSendMailUseCase(sentMailRepo, threadGroupRepo, senderRepo, discordClient)
	getSettingsUC := settingsuc.NewGetUserSettingsUseCase(userSettingRepo)
	updateSettingsUC := settingsuc.NewUpdateUserSettingsUseCase(userSettingRepo)

	// Handlers
	mailHandler := handler.NewMailHandler(getMailsUC, updateStateUC, deleteMailUC, syncMailsUC, userSettingRepo, domainRepo)
	threadHandler := handler.NewThreadHandler(getThreadUC, userSettingRepo, domainRepo)
	sendHandler := handler.NewSendHandler(sendMailUC)
	settingsHandler := handler.NewSettingsHandler(getSettingsUC, updateSettingsUC)
	domainHandler := handler.NewDomainHandler(domainRepo)
	systemSettingHandler := handler.NewSystemSettingHandler(systemSettingRepo)

	// Authenticated routes
	api := e.Group("/api", middleware.FirebaseAuth(fbAuth, userRepo))

	// Mail routes
	api.GET("/mails", mailHandler.GetMails)
	api.GET("/mails/:s3Key", mailHandler.GetMail)
	api.PATCH("/mails/:s3Key/read", mailHandler.UpdateReadStatus)
	api.PATCH("/mails/:s3Key/star", mailHandler.UpdateStarStatus)
	api.DELETE("/mails/:s3Key", mailHandler.DeleteMail)
	api.POST("/mails/sync", mailHandler.SyncMails)
	api.GET("/mails/recipients", mailHandler.GetRecipients)

	// Thread routes
	api.GET("/threads", threadHandler.ListThreads)
	api.GET("/threads/:threadId", threadHandler.GetThread)

	// Send routes
	api.POST("/send", sendHandler.SendMail)

	// Settings routes
	api.GET("/settings", settingsHandler.GetSettings)
	api.PUT("/settings", settingsHandler.UpdateSettings)

	// Domain routes
	api.GET("/domains", domainHandler.ListDomains)
	api.POST("/domains", domainHandler.CreateDomain)
	api.PUT("/domains/:id", domainHandler.UpdateDomain)
	api.DELETE("/domains/:id", domainHandler.DeleteDomain)

	// System settings (admin only)
	api.GET("/system/settings", systemSettingHandler.Get)
	api.PUT("/system/settings", systemSettingHandler.Update)

	return e
}
