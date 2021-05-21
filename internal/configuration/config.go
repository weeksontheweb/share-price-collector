package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"share-price-collector/internal/database"
	settingsfile "share-price-collector/internal/settings-file"
)

//The share struct includes poll details for the
//individual shares (overrrides defaults)
type share struct {
	Code         string `json:"code"`
	Description  string `json:"description"`
	PollStart    string `json:"startPoll"`
	PollEnd      string `json:"endPoll"`
	PollInterval int    `json:"interval"`
}

type ConfigDetails struct {
	PollStart    string  `json:"startPoll"`
	PollEnd      string  `json:"endPoll"`
	PollInterval int     `json:"interval"`
	Shares       []share `json:"shareCodes"`
}

func (c ConfigDetails) ReadConfig(i interface{}) ConfigDetails {

	var newConfig ConfigDetails
	var err error

	fmt.Printf("type = %s\n", reflect.TypeOf(i).String())

	reflect.ValueOf(i)

	switch bb := i.(type) {
	case *os.File:
		fmt.Println("In case for *os.File.")

		newConfig, err = readConfigFromFile()

	case database.SharesDB:
		fmt.Println("In case for sql.DB.")

		newConfig, err = readConfigFromDB(bb)

	default:
		log.Fatal()
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Returned PollStart = %s\n", newConfig.PollStart)

	return newConfig
}

func readConfigFromFile() (ConfigDetails, error) {

	jsonFile, err := settingsfile.LoadSettingsFile()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened settings.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var sharesConfig ConfigDetails

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'sharesConfig' which we defined above
	err = json.Unmarshal(byteValue, &sharesConfig)

	if err != nil {
		fmt.Printf("Unmarshal error = %s\n", err)
	}

	return sharesConfig, nil
}

func readConfigFromDB(db database.SharesDB) (ConfigDetails, error) {

	var sharesConfig ConfigDetails

	//Load defaults (store these in database at some point)
	sharesConfig.PollStart = "08:00"
	sharesConfig.PollEnd = "17:00"
	sharesConfig.PollInterval = 5

	retrievedRecords := db.RetrieveShares()

	for _, shareRecord := range retrievedRecords {
		var shareRecordConfig share

		shareRecordConfig.Code = shareRecord.ShareCode
		shareRecordConfig.Description = shareRecord.ShareDescription
		shareRecordConfig.PollStart = shareRecord.PollStart
		shareRecordConfig.PollEnd = shareRecord.PollEnd
		shareRecordConfig.PollInterval = shareRecord.PollInterval

		sharesConfig.Shares = append(sharesConfig.Shares, shareRecordConfig)

	}

	return sharesConfig, nil

}
