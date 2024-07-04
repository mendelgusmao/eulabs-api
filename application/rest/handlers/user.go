package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/mendelgusmao/eulabs-api/application/rest"
	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
)

type UserService interface {
	Authorize(context.Context, dto.UserCredentials) (*dto.User, error)
}

type UserHandler struct {
	service   UserService
	jwtConfig rest.JWTConfig
}

func NewUserHandler(e *echo.Echo, jwtConfig rest.JWTConfig, s UserService) {
	h := &UserHandler{
		service:   s,
		jwtConfig: jwtConfig,
	}

	e.POST("/users/authenticate", h.authenticate)
}

func (h *UserHandler) authenticate(c echo.Context) error {
	ctx := c.Request().Context()
	credentials := dto.UserCredentials{}

	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rest.Error(err))
	}

	user, err := h.service.Authorize(ctx, credentials)

	if err != nil {
		if err == domain.ErrCredentialsDontMatch {
			c.NoContent(http.StatusUnauthorized)
		}

		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	claims := &rest.JWTClaims{
		ID:    user.ID,
		Name:  user.Name,
		Admin: user.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.jwtConfig.Expiration)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(h.jwtConfig.Secret)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
