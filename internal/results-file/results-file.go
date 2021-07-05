package resultsfile

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func AppendToResultsFile(code string, midPrice float64, bidPrice float64, askPrice float64) {

	//var file *os.File

	_, err := os.Stat("../../share-poll-results.csv")

	if os.IsNotExist(err) {
		os.Create("../../share-poll-results.csv")
		//log.Fatal("File does not exist.")
	}

	Myfile, err := os.OpenFile("../../share-poll-results.csv", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Unable to open file")
	}

	//len, err := Myfile.WriteString(" Hello India")
	strconv.FormatFloat(midPrice, 'f', 6, 64)
	t := time.Now()

	//len, err := Myfile.WriteString(code + ",")
	len, err := Myfile.WriteString(strings.Trim(code, " ") + "," + t.Format("2006-01-02 15:04:05") + "," + strconv.FormatFloat(midPrice, 'f', 6, 64) + "," + strconv.FormatFloat(bidPrice, 'f', 6, 64) + "," + strconv.FormatFloat(askPrice, 'f', 6, 64) + "\n")

	if err != nil {
		fmt.Printf("Error = %s\n", err)
	}

	if len == 0 {
		fmt.Printf("File is opened in readonly mode")
	} else {
		fmt.Printf("\n%d characters written into file", len)
	}
	Myfile.Close()

}
