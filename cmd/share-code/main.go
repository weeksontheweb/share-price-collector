package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	config "share-price-collector/internal/configuration"
	"share-price-collector/internal/database"
	settingsfile "share-price-collector/internal/settings-file"

	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

var (
	host         string
	port         int
	user         string
	passwd       string
	dbname       string
	nodb         bool
	output       string
	sharecode    string
	sharedesc    string
	pollstart    string
	pollend      string
	pollinterval int
)

func main() {

	app := &cli.App{
		Name:  "share-code",
		Usage: "Update/delete or list config share codes",

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
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a share code to the config.",

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
					&cli.StringFlag{
						Name:        "sharecode",
						Value:       "",
						Usage:       "Unique share code.",
						Destination: &sharecode,
					},
					&cli.StringFlag{
						Name:        "sharedesc",
						Value:       "",
						Usage:       "Share description.",
						Destination: &sharedesc,
					},
					&cli.StringFlag{
						Name:        "pollstart",
						Value:       "",
						Usage:       "Polling start time.",
						Destination: &pollstart,
					},
					&cli.StringFlag{
						Name:        "pollend",
						Value:       "",
						Usage:       "Polling end time.",
						Destination: &pollend,
					},
					&cli.IntFlag{
						Name:        "pollinterval",
						Value:       0,
						Usage:       "Polling interval (minutes).",
						Destination: &pollinterval,
					},
				},
			},
			{
				Name:  "remove",
				Usage: "Remove a share code from the config.",

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
					&cli.StringFlag{
						Name:        "sharecode",
						Value:       "",
						Usage:       "Unique share code.",
						Destination: &sharecode,
					},
					&cli.StringFlag{
						Name:        "sharedesc",
						Value:       "",
						Usage:       "Share description.",
						Destination: &sharedesc,
					},
					&cli.StringFlag{
						Name:        "pollstart",
						Value:       "",
						Usage:       "Polling start time.",
						Destination: &pollstart,
					},
					&cli.StringFlag{
						Name:        "pollend",
						Value:       "",
						Usage:       "Polling end time.",
						Destination: &pollend,
					},
					&cli.IntFlag{
						Name:        "pollinterval",
						Value:       0,
						Usage:       "Polling interval (minutes).",
						Destination: &pollinterval,
					},
				},
			},
			{
				Name:  "list",
				Usage: "List all share codes in the config.",
				/*				Action: func(c *cli.Context) error {
								ns, nsErr := net.LookupIP(c.String("host"))
								if nsErr != nil {
									return nsErr
								}
								for _, v := range ns {
									fmt.Println(v)
								}
								return nil
							},*/
			},
		},
	}

	app.Action = actionShareCode

	app.Commands[0].Action = shareCodeAdd
	app.Commands[1].Action = shareCodeRemove
	app.Commands[2].Action = shareCodeList

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

//Anything called by the cli library has to have the
//signature func (c *cli.Context) error
func actionShareCode(c *cli.Context) error {

	//var myConfig config.ConfigDetails

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

			fmt.Println("database requested.")
			//Open the requested database.
			//			var db database.SharesDB

			//			db, err = db.ConnectToDatabase(host, port, user, passwd, dbname)
		} else {
			fmt.Println("database not requested.")
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

func shareCodeAdd(c *cli.Context) error {

	fmt.Println("Share code add.")

	fmt.Printf("host = %s\n", host)
	fmt.Printf("port = %d\n", port)
	fmt.Printf("user = %s\n", user)
	fmt.Printf("passwd = %s\n", passwd)
	fmt.Printf("dbname = %s\n", dbname)
	fmt.Printf("sharecode = %s\n", sharecode)
	fmt.Printf("sharedesc = %s\n", sharedesc)

	var myConfig config.ConfigDetails

	//See if the database is requested in the command line.
	databaseRequested, err := requireToUseDatabase()

	if err != nil {
		panic(err)
	} else {
		if databaseRequested {

			fmt.Println("database requested.")

			var db database.SharesDB

			db, err = db.ConnectToDatabase(host, port, user, passwd, dbname)

			fmt.Println("aaaa")

			if err != nil {
				fmt.Println("bbbb")
				panic(err)
			}

			myConfig.AddShareToConfig(db, sharecode, sharedesc, pollstart, pollend, pollinterval)

		} else {
			fmt.Println("database not requested.")

			jsonFile, err := settingsfile.LoadSettingsFile()

			if err != nil {
				log.Fatal(err)
			}

			myConfig.AddShareToConfig(jsonFile, sharecode, sharedesc, pollstart, pollend, pollinterval)
		}
	}

	return nil
}

func shareCodeRemove(c *cli.Context) error {

	fmt.Println("Share code remove.")

	var myConfig config.ConfigDetails

	//See if the database is requested in the command line.
	databaseRequested, err := requireToUseDatabase()

	if err != nil {
		panic(err)
	} else {
		if databaseRequested {

			fmt.Println("database requested.")

			var db database.SharesDB

			db, err = db.ConnectToDatabase(host, port, user, passwd, dbname)

			fmt.Println("aaaa")

			if err != nil {
				panic(err)
			}

			myConfig.RemoveShareFromConfig(db, sharecode)

		} else {
			fmt.Println("database not requested.")

			jsonFile, err := settingsfile.LoadSettingsFile()

			if err != nil {
				log.Fatal(err)
			}

			myConfig.RemoveShareFromConfig(jsonFile, sharecode)
		}
	}
	return nil
}

func shareCodeList(c *cli.Context) error {

	fmt.Println("Share code list.")
	//See if the database is requested in the command line.
	databaseRequested, err := requireToUseDatabase()

	if err != nil {
		panic(err)
	} else {
		if databaseRequested {

			fmt.Println("database requested.")

			var db database.SharesDB

			db, err = db.ConnectToDatabase(host, port, user, passwd, dbname)

			fmt.Println("aaaa")

			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("database not requested.")
		}
	}

	return nil
}
