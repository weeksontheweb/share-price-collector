package config

import (
	"encoding/json"
	"fmt"
	"io"
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

func writeConfigToFile(config ConfigDetails) error {

	err := settingsfile.BackupSettingsFile()

	if err != nil {
		return nil
	}

	file, _ := json.MarshalIndent(config, "", " ")

	err = settingsfile.CreateSettingsFile(file)

	if err != nil {
		return nil
	}

	return nil
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

func copyFile() {

	// Open original file
	original, err := os.Open("original.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	// Create new file
	new, err := os.Create("new.txt")

	if err != nil {
		log.Fatal(err)

	}
	defer new.Close()

	//This will copy
	bytesWritten, err := io.Copy(new, original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes Written: %d\n", bytesWritten)
}

func (c ConfigDetails) AddShareToConfig(i interface{}, shareCode string, shareDescription string, pollStart string, pollEnd string, pollInterval int) error {

	//var newConfig ConfigDetails
	//var err error

	//switch bb := i.(type) {
	switch bb := i.(type) {
	case *os.File:
		fmt.Println("In case for *os.File. shareCode = " + shareCode)

		newConfig, err := readConfigFromFile()

		if err != nil {
			log.Panic(err)
		}

		//Add share to struct and marshall.
		var newShare share

		newShare.Code = shareCode
		newShare.Description = shareDescription
		newShare.PollStart = pollStart
		newShare.PollEnd = pollEnd
		newShare.PollInterval = pollInterval

		newConfig.Shares = append(newConfig.Shares, newShare)

		writeConfigToFile(newConfig)

	case database.SharesDB:
		fmt.Println("In case for sql.DB. shareCode = " + shareCode)

		//Don't need to read the config from the db.
		//Just insert  a new record.

		bb.AddShareCode(shareCode, shareDescription, pollStart, pollEnd, pollInterval)

	default:
		log.Fatal()
	}

	return nil
}
