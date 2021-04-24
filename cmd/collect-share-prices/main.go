package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	fmt.Printf("host = %s\n", host)
	fmt.Printf("port = %d\n", port)
	fmt.Printf("user = %s\n", user)
	fmt.Printf("passwd = %s\n", passwd)
	fmt.Printf("dbname = %s\n", dbname)

	//readConfig()

	midPrice, bidPrice, askPrice, err := scrapeSharePrices("TEP")
	fmt.Printf("TEP\t\t%f\t%f\t%f\n", midPrice, bidPrice, askPrice)

	if err != nil {
		panic(err)
	}

	midPrice, bidPrice, askPrice, err = scrapeSharePrices("ARCM")
	fmt.Printf("ARCM\t\t%f\t%f\t%f\n", midPrice, bidPrice, askPrice)

	if err != nil {
		panic(err)
	}

	midPrice, bidPrice, askPrice, err = scrapeSharePrices("XTR")
	fmt.Printf("XTR\t\t%f\t%f\t%f\n", midPrice, bidPrice, askPrice)

	if err != nil {
		panic(err)
	}

	midPrice, bidPrice, askPrice, err = scrapeSharePrices("BIRG")
	fmt.Printf("BIRG\t\t%f\t%f\t%f\n", midPrice, bidPrice, askPrice)

	if err != nil {
		panic(err)
	}

	/*
		//See if the database is requested in the command line.
		databaseRequested, err := useDatabase()

		if err != nil {
			panic(err)
		} else {
			if databaseRequested {
				db, err := ConnectToDatabase(host, port, user, passwd, dbname)

				if err != nil {
					panic(err)
				}
			} else {

			}
		}
	*/

	//fmt.Printf("Here %d\n", db.Stats().WaitDuration)
	return nil
}

func ConnectToDatabase(host string, port int, user string, password string, dbname string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Println("Got an error here")
		return nil, err
	}

	return newdb, nil
}

func scrapeSharePrices(shareCode string) (float64, float64, float64, error) {

	var midPrice float64
	var bidPrice float64
	var askPrice float64

	resp, err := http.Get("https://www.lse.co.uk/shareprice.asp?shareprice=" + shareCode)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if bytes.Contains(body, []byte("html")) {
		webBody := string(body[:])

		lookFor := "<span data-item=\"" + shareCode + ".L\" data-field=\"MID_PRICE\"  data-flash=\"true\">"
		midPrice = float64(scrapeSharePrice(webBody, lookFor))

		lookFor = "<span data-item=\"" + shareCode + ".L\" data-field=\"BID\"  data-flash=\"true\">"
		bidPrice = float64(scrapeSharePrice(webBody, lookFor))

		lookFor = "<span data-item=\"" + shareCode + ".L\" data-field=\"ASK\"  data-flash=\"true\">"
		askPrice = float64(scrapeSharePrice(webBody, lookFor))
	}

	return midPrice, bidPrice, askPrice, nil
}

func scrapeSharePrice(body string, searchFor string) float64 {

	index := strings.Index(body, searchFor)
	elementLength := len(searchFor)
	truncatedBody := body[index+elementLength:]
	nextElementPosition := strings.Index(truncatedBody, "<")

	fValue := strings.Replace(truncatedBody[:nextElementPosition], ",", "", -1)

	s, err := strconv.ParseFloat(fValue, 64)

	if err != nil {
		fmt.Printf("%T, %v\n", s, s)
	}

	return s
}

/*
//Determines whether a database is requested on the command line.
// If no database flags are requested then no database requested.
//Error if partial database details given
func useDatabase() (bool, error) {

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
*/
