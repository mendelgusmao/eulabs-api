package application

import (
	"log"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mendelgusmao/eulabs-api/application/rest"
	"github.com/mendelgusmao/eulabs-api/infrastructure/database"
	"github.com/mendelgusmao/eulabs-api/repository"
	"github.com/mendelgusmao/eulabs-api/service"
)

type ServerConfig struct {
	Address string
	DSN     string
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

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	e.Validator = rest.NewValidator()

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	rest.NewProductHandler(e, productService)

	return &Server{
		config: config,
		echo:   e,
	}
}

func (s *Server) Serve() error {
	return s.echo.Start(s.config.Address)
}
