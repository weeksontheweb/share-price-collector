package settingsfile

import (
	"fmt"
	"log"
	"os"
)

func LoadSettingsFile() (*os.File, error) {

	var file *os.File

	_, err := os.Stat("../../settings.json")

	if os.IsNotExist(err) {
		log.Fatal("File does not exist.")
	}

	file, err = os.OpenFile("../../settings.json", os.O_RDWR, 0644)

	if err != nil {
		fmt.Println(err)
	}

	return file, err

}
