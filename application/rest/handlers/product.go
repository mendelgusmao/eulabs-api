package handlers

import (
	"context"
	"net/http"
	"strconv"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/mendelgusmao/eulabs-api/application/rest"
	"github.com/mendelgusmao/eulabs-api/application/rest/middleware"
	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
)

type ProductService interface {
	GetMany(context.Context) ([]dto.Product, error)
	GetOne(context.Context, int64) (*dto.Product, error)
	Create(context.Context, dto.BaseProduct) (*dto.Product, error)
	Update(context.Context, dto.UpdateProduct) (*dto.Product, error)
	Delete(context.Context, int64) error
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(e *echo.Group, echoJWTConfig echojwt.Config, s ProductService) {
	h := &ProductHandler{
		service: s,
	}

	g := e.Group("/products")
	g.Use(echojwt.WithConfig(echoJWTConfig))

	g.GET("", h.list)
	g.GET("/:id", h.get)
	g.POST("", h.create, middleware.AdminMiddleware)
	g.PUT("/:id", h.update, middleware.AdminMiddleware)
	g.DELETE("/:id", h.delete, middleware.AdminMiddleware)
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
	createdProduct, err := h.service.Create(ctx, product)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	return c.JSON(http.StatusCreated, createdProduct)
}

func (h *ProductHandler) update(c echo.Context) error {
	ctx := c.Request().Context()

	var (
		product   = dto.UpdateProduct{}
		productId int
	)

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rest.Error(err))
	}

	if id := c.Param("id"); id != "" {
		var err error
		productId, err = strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, rest.ErrInvalidIdType)
		}
	} else {
		return c.JSON(http.StatusBadRequest, rest.ErrEmptyId)
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, rest.Error(err))
	}

	product.ID = int64(productId)
	updatedProduct, err := h.service.Update(ctx, product)

	if err != nil {
		if err == domain.ErrNotFound {
			return c.JSON(http.StatusNotFound, rest.Error(err))
		}

		return c.JSON(http.StatusInternalServerError, rest.Error(err))
	}

	return c.JSON(http.StatusOK, updatedProduct)
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
