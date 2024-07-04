package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/mendelgusmao/eulabs-api/application/rest"
	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
)

type ProductService interface {
	GetMany(context.Context) ([]dto.Product, error)
	GetOne(context.Context, int64) (*dto.Product, error)
	Create(context.Context, dto.BaseProduct) error
	Update(context.Context, dto.Product) error
	Delete(context.Context, int64) error
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(e *echo.Echo, jwtConfig rest.JWTConfig, s ProductService) {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return rest.JWTClaims{}
		},
		SigningKey: jwtConfig.Secret,
	}

	h := &ProductHandler{
		service: s,
	}

	e.GET("/products", h.list)
	e.GET("/products/:id", h.get)
	e.POST("/products", h.create, echojwt.WithConfig(config))
	e.PUT("/products/:id", h.update, echojwt.WithConfig(config))
	e.DELETE("/products/:id", h.delete, echojwt.WithConfig(config))
}

func (h *ProductHandler) list(c echo.Context) error {
	ctx := c.Request().Context()
	products, err := h.service.GetMany(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	if len(products) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) get(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, rest.ErrInvalidIdType)
	}

	productId := int64(id)
	product, err := h.service.GetOne(ctx, productId)

	if err != nil {
		if err == domain.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) create(c echo.Context) error {
	product := dto.BaseProduct{}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rest.Error(err))
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rest.Error(err))
	}

	ctx := c.Request().Context()

	if err := h.service.Create(ctx, product); err != nil {
		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) update(c echo.Context) error {
	var product = dto.Product{}

	if id := c.Param("id"); id != "" {
		productId, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, rest.ErrInvalidIdType)
		}

		product.ID = int64(productId)
	} else {
		return c.JSON(http.StatusBadRequest, rest.ErrEmptyId)
	}

	ctx := c.Request().Context()

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rest.Error(err))
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusBadRequest, rest.Error(err))
	}

	if err := h.service.Update(ctx, product); err != nil {
		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) delete(c echo.Context) error {
	var productId int

	if id := c.Param("id"); id != "" {
		var err error
		productId, err = strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, rest.ErrInvalidIdType)
		}
	} else {
		return c.JSON(http.StatusBadRequest, rest.ErrEmptyId)
	}

	ctx := c.Request().Context()

	if err := h.service.Delete(ctx, int64(productId)); err != nil {
		if err == domain.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	return c.NoContent(http.StatusNoContent)
}
