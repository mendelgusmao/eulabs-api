package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ConfigSpecification struct {
	Address       string `default:":8080"`
	DSN           string
	JWTSecret     []byte        `envconfig:"jwt_secret"`
	JWTExpiration time.Duration `envconfig:"jwt_expiration"`
}

var config ConfigSpecification

func init() {
	if err := envconfig.Process("EULABS_API", &config); err != nil {
		log.Fatal(err.Error())
	}
}
