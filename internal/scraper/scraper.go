package scraper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func RetrieveSharePrices(shareCode string) (float64, float64, float64, error) {

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
		midPrice = float64(retrieveSharePrice(webBody, lookFor))

		lookFor = "<span data-item=\"" + shareCode + ".L\" data-field=\"BID\"  data-flash=\"true\">"
		bidPrice = float64(retrieveSharePrice(webBody, lookFor))

		lookFor = "<span data-item=\"" + shareCode + ".L\" data-field=\"ASK\"  data-flash=\"true\">"
		askPrice = float64(retrieveSharePrice(webBody, lookFor))
	}

	return midPrice, bidPrice, askPrice, nil
}

func retrieveSharePrice(body string, searchFor string) float64 {

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
