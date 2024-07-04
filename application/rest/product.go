package rest

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
)

type ResponseError struct {
	Message string `json:"message"`
}

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

func NewProductHandler(e *echo.Echo, s ProductService) {
	h := &ProductHandler{
		service: s,
	}

	e.GET("/products", h.list)
	e.GET("/products/:id", h.get)
	e.POST("/products", h.create)
	e.PUT("/products/:id", h.update)
	e.DELETE("/products/:id", h.delete)
}

func (h *ProductHandler) list(c echo.Context) error {
	ctx := c.Request().Context()
	products, err := h.service.GetMany(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error(err))
	}

	if len(products) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) get(c echo.Context) error {
	log.Println("get one")
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidIdType)
	}

	productId := int64(id)
	product, err := h.service.GetOne(ctx, productId)

	if err != nil {
		if err == domain.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, Error(err))
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) create(c echo.Context) error {
	product := dto.BaseProduct{}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Error(err))
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Error(err))
	}

	ctx := c.Request().Context()

	if err := h.service.Create(ctx, product); err != nil {
		return c.JSON(http.StatusInternalServerError, Error(err))
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) update(c echo.Context) error {
	var product = dto.Product{}

	if id := c.Param("id"); id != "" {
		productId, err := strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrInvalidIdType)
		}

		product.ID = int64(productId)
	} else {
		return c.JSON(http.StatusBadRequest, ErrEmptyId)
	}

	ctx := c.Request().Context()

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Error(err))
	}

	if err := c.Validate(product); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	if err := h.service.Update(ctx, product); err != nil {
		return c.JSON(http.StatusInternalServerError, Error(err))
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) delete(c echo.Context) error {
	var productId int

	if id := c.Param("id"); id != "" {
		var err error
		productId, err = strconv.Atoi(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrInvalidIdType)
		}
	} else {
		return c.JSON(http.StatusBadRequest, ErrEmptyId)
	}

	ctx := c.Request().Context()

	if err := h.service.Delete(ctx, int64(productId)); err != nil {
		if err == domain.ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, Error(err))
	}

	return c.NoContent(http.StatusNoContent)
}
