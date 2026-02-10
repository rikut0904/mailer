package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
	fbinfra "github.com/rikut0904/mailer-backend/internal/infrastructure/firebase"
	"gorm.io/gorm"
)

func FirebaseAuth(fbAuth *fbinfra.FirebaseAuth, userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header"})
			}

			token, err := fbAuth.VerifyIDToken(c.Request().Context(), parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			c.Set("uid", token.UID)

			role := "user"
			user, err := userRepo.GetByUID(token.UID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					_ = userRepo.Upsert(&entity.User{
						UID:  token.UID,
						Role: "user",
					})
				} else {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load user"})
				}
			} else if user.Role != "" {
				role = user.Role
			}

			c.Set("role", role)
			return next(c)
		}
	}
}
