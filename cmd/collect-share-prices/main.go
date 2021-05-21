package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	config "share-price-collector/internal/configuration"
	"share-price-collector/internal/database"
	"share-price-collector/internal/scraper"
	settingsfile "share-price-collector/internal/settings-file"

	//"github.com/weeksontheweb/share-price-collector/internal/database"
	//"github.com/weeksontheweb/share-price-collector/internal/database"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

var (
	language string
	host     string
	port     int
	user     string
	passwd   string
	dbname   string
	nodb     bool
	output   string
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
			&cli.IntFlag{
				Name:        "port",
				Value:       0,
				Usage:       "Database port",
				Destination: &port,
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
				Required:    false,
				Usage:       "Database username",
				Destination: &user,
			},
			&cli.StringFlag{
				Name:        "passwd",
				Value:       "",
				Usage:       "Database password",
				Destination: &passwd,
			},
			&cli.BoolFlag{
				Name:        "nodb",
				Value:       false,
				Required:    false,
				Usage:       "Use a backend database",
				Destination: &nodb,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "screen",
				Usage:       "Results output type",
				Destination: &output,
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

	var myConfig config.ConfigDetails

	fmt.Printf("host = %s\n", host)
	fmt.Printf("port = %d\n", port)
	fmt.Printf("user = %s\n", user)
	fmt.Printf("passwd = %s\n", passwd)
	fmt.Printf("dbname = %s\n", dbname)

	//See if the database is requested in the command line.
	databaseRequested, err := requireToUseDatabase()

	if err != nil {
		panic(err)
	} else {
		if databaseRequested {

			//Open the requested database.
			var db database.SharesDB

			db, err = db.ConnectToDatabase(host, port, user, passwd, dbname)

			fmt.Println("aaaa")

			if err != nil {
				panic(err)
			}

			myConfig = myConfig.ReadConfig(db)

			for _, share := range myConfig.Shares {
				fmt.Printf("After ReadConfig. Share Code = %s\n", share.Code)
				fmt.Printf("After ReadConfig. Share Description = %s\n", share.Description)
				fmt.Printf("After ReadConfig. Share PollStart = %s\n", share.PollStart)
				fmt.Printf("After ReadConfig. Share PollEnd = %s\n", share.PollEnd)
				fmt.Printf("After ReadConfig. Share PollInterval = %d\n", share.PollInterval)
				midPrice, bidPrice, askPrice, err := scraper.RetrieveSharePrices(share.Code)
				fmt.Printf("%s\t\t%f\t%f\t%f\n", share.Code, midPrice, bidPrice, askPrice)

				if err != nil {
					panic(err)
				}
			}

		} else {
			jsonFile, err := settingsfile.LoadSettingsFile()

			if err != nil {
				log.Fatal(err)
			}

			myConfig = myConfig.ReadConfig(jsonFile)

			for _, share := range myConfig.Shares {
				fmt.Printf("After ReadConfig. Share Code = %s\n", share.Code)
				midPrice, bidPrice, askPrice, err := scraper.RetrieveSharePrices(share.Code)
				fmt.Printf("%s\t\t%f\t%f\t%f\n", share.Code, midPrice, bidPrice, askPrice)

				if err != nil {
					panic(err)
				}
			}
		}
	}

	return nil
}

//Determines whether a database is requested on the command line.
// If no database flags are requested then no database requested.
//Error if partial database details given
func requireToUseDatabase() (bool, error) {

	var nCount int

	if host != "" {
		nCount++
	}
	if user != "" {
		nCount++
	}
	if passwd != "" {
		nCount++
	}

	switch nCount {
	case 0:
		return false, nil
	case 3:
		return true, nil
	default:
		return false, errors.New("an error")
	}
}
