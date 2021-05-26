package settingsfile

import (
	"fmt"
	"io/ioutil"
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

func BackupSettingsFile() error {

	_, err := os.Stat("../../settings.json.backup")

	//Remove the curtrent backup if it exists.
	if !os.IsNotExist(err) {
		fmt.Println("Backup file exists.")

		err := os.Remove("../../settings.json.backup")

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	err = os.Rename("../../settings.json", "../../settings.json.backup")

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func CreateSettingsFile(file []byte) error {

	err := ioutil.WriteFile("../../settings.json", file, 0644)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
