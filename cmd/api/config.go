package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type ConfigSpecification struct {
	Address   string `default:":8080"`
	DSN       string
	JWTSecret string `envconfig:"jwt_secret"`
}

var config ConfigSpecification

func init() {
	if err := envconfig.Process("EULABS_API", &config); err != nil {
		log.Fatal(err.Error())
	}
}
