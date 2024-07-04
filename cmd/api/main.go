package main

import (
	"log"
	"os"

	"github.com/mendelgusmao/eulabs-api/application"
	"github.com/mendelgusmao/eulabs-api/infrastructure/database"
	"github.com/urfave/cli"
)

var app = &cli.App{
	Action: func(cCtx *cli.Context) error {
		action := ""

		if cCtx.NArg() > 0 {
			action = cCtx.Args().Get(0)
		}

		if action == "" {
			serverConfig := application.ServerConfig{
				Address:       config.Address,
				DSN:           config.DSN,
				JWTSecret:     config.JWTSecret,
				JWTExpiration: config.JWTExpiration,
			}

			server := application.NewServer(serverConfig)
			log.Fatal(server.Serve())
		} else if action == "migratedb" {
			log.Println("migrating database")

			if err := database.Migrate(config.DSN); err != nil {
				log.Fatal(err)
			}

			log.Println("database migration finished")
		} else if action == "populatedb" {
			log.Println("populating database")

			if err := database.Populate(config.DSN); err != nil {
				log.Fatal(err)
			}

			log.Println("database population finished")
		} else {
			log.Fatal("unrecognized action:", action)
		}

		return nil
	},
}

// @title           Eulabs Products API
// @version         1.0
// @description     A CRUD API to deal with products.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Developer
// @contact.email  mendelson.gusmao@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      https://ludovic.fawn-beaver.ts.net:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer JWT

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
