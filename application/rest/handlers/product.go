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
	Create(context.Context, dto.BaseProduct) error
	Update(context.Context, dto.Product) error
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
	g.GET(":id", h.get)
	g.POST("", h.create, middleware.AdminMiddleware)
	g.PUT(":id", h.update, middleware.AdminMiddleware)
	g.DELETE(":id", h.delete, middleware.AdminMiddleware)
}

// @Summary List all products
// @Description Returns all registered products
// @Tags products
// @Produce json
// @Success 200 {array} dto.Product
// @Success 204 {string} empty
// @Failure 500 {object} rest.Error
// @Router /products [get]
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

// @Summary Get a product by ID
// @Description Returns details of a specific product
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.Product
// @Success 404 {string} string "Not Found"
// @Failure 400 {object} rest.Error
// @Failure 500 {object} rest.Error
// @Router /products/{id} [get]
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

// @Summary Create a new product
// @Description Adds a new product to the system
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.BaseProduct true "Product to be created"
// @Success 201 {object} dto.Product
// @Failure 422 {object} rest.Error
// @Failure 500 {object} rest.Error
// @Router /products [post]
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

// @Summary Update a product
// @Description Updates the details of an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.Product true "Product to be updated"
// @Success 200 {object} dto.Product
// @Failure 400 {object} rest.Error
// @Failure 422 {object} rest.Error
// @Failure 500 {object} rest.Error
// @Router /products/{id} [put]
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

// @Summary Delete a product
// @Description Removes a product from the system
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} rest.Error
// @Failure 404 {string} string "Not Found"
// @Failure 500 {object} rest.Error
// @Router /products/{id} [delete]
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
