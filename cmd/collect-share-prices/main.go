package main

import (
	"fmt"
	"log"
	"os"

	//"github.com/weeksontheweb/share-price-collector/internal/database"
	//"github.com/weeksontheweb/share-price-collector/internal/database"

	"github.com/urfave/cli/v2"
)

var (
	language string
	host     string
	user     string
	passwd   string
	dbname   string
)

func main() {

	app := &cli.App{
		Name:  "share-price-grabber",
		Usage: "Grabs share prices from lse.co.uk",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       "",
				Usage:       "Host where database is located",
				Destination: &host,
			},
			&cli.StringFlag{
				Name:        "dbname",
				Value:       "",
				Usage:       "Database name",
				Destination: &dbname,
			},
			&cli.StringFlag{
				Name:        "user",
				Value:       "",
				Required:    true,
				Usage:       "Database username",
				Destination: &user,
			},
			&cli.StringFlag{
				Name:        "passwd",
				Value:       "",
				Usage:       "Database password",
				Destination: &passwd,
			},
		},
	}

	app.Action = actionShareGrabber

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

//Anything called by the cli library has to have the
//signature func (c *cli.Context) error
func actionShareGrabber(c *cli.Context) error {
	//var db database.ShareDB

	fmt.Printf("User = %s\n", host)
	fmt.Printf("User = %s\n", user)
	fmt.Printf("User = %s\n", passwd)
	fmt.Printf("User = %s\n", dbname)

	//db.ConnectToDatabase(host, 5432, user, passwd, dbname)

	return nil
}
