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

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
