package application

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	handlers.NewProductHandler(e, jwtConfig, productService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	handlers.NewUserHandler(e, jwtConfig, userService)

	return &Server{
		config: config,
		echo:   e,
	}
}

func (s *Server) Serve() error {
	return s.echo.Start(s.config.Address)
}
