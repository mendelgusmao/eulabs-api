package application

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/mendelgusmao/eulabs-api/cmd/api/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/mendelgusmao/eulabs-api/application/rest"
	"github.com/mendelgusmao/eulabs-api/application/rest/handlers"
	"github.com/mendelgusmao/eulabs-api/infrastructure/database"
	"github.com/mendelgusmao/eulabs-api/repository"
	"github.com/mendelgusmao/eulabs-api/service"
)

type ServerConfig struct {
	Address       string
	DSN           string
	JWTSecret     []byte
	JWTExpiration time.Duration
}

type Server struct {
	config ServerConfig
	echo   *echo.Echo
}

func NewServer(config ServerConfig) *Server {
	db, err := database.Connect(config.DSN)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	jwtConfig := rest.JWTConfig{
		Secret:     config.JWTSecret,
		Expiration: config.JWTExpiration,
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	e.Validator = rest.NewValidator()
	echoJWTConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &rest.JWTClaims{}
		},
		SigningKey: jwtConfig.Secret,
		ErrorHandler: func(c echo.Context, err error) error {
			c.JSON(http.StatusBadRequest, rest.Error(err))
			return nil
		},
	}

	g := e.Group("/api/v1")

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	handlers.NewProductHandler(g, echoJWTConfig, productService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	handlers.NewUserHandler(g, jwtConfig, userService)

	return &Server{
		config: config,
		echo:   e,
	}
}

func (s *Server) Serve() error {
	return s.echo.Start(s.config.Address)
}
