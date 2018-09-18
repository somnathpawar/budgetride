package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"os"

)

func UberCabs(startLn, startLg, endLn, endLg string) (CabList, error) {
	type cabPrices struct {
		LocalizedDisplayName string  `json:"localized_display_name"`
		Distance             float64 `json:"distance"`
		DisplayName          string  `json:"display_name"`
		ProductID            string  `json:"product_id"`
		HighEstimate         float64 `json:"high_estimate"`
		LowEstimate          float64 `json:"low_estimate"`
		Duration             int     `json:"duration"`
		Estimate             string  `json:"estimate"`
		CurrencyCode         string  `json:"currency_code"`
	}
	type estimatePrices struct {
		Prices []cabPrices `json:"prices"`
	}

	priceBody := estimatePrices{}
	uberPrices, err := uberRequest("https://api.uber.com/v1.2/estimates/price?start_latitude=" + startLn + "&start_longitude=" + startLg + "&end_latitude=" + endLn + "&end_longitude=" + endLg)
	err = json.Unmarshal(uberPrices, &priceBody)
	if err != nil {
		return nil, err
	}

	type cabTimes struct {
		LocalizedDisplayName string `json:"localized_display_name"`
		Estimate             int    `json:"estimate"`
		DisplayName          string `json:"display_name"`
		ProductID            string `json:"product_id"`
	}
	type estimateTimes struct {
		Times []cabTimes `json:"times"`
	}
	estimateBody := estimateTimes{}
	uberTimes, err := uberRequest("https://api.uber.com/v1.2/estimates/time?start_latitude=" + startLn + "&start_longitude=" + startLg)
	err = json.Unmarshal(uberTimes, &estimateBody)
	if err != nil {
		return nil, err
	}

	uberCabs := CabList{}
	for _, cab := range priceBody.Prices {
		uberCab := CabResponse{}
		uberCab.Company = "Uber"
		uberCab.Cab = cab.DisplayName
		uberCab.Estimate = cab.LowEstimate
		for _, cabTimes := range estimateBody.Times {
			if cabTimes.ProductID == cab.ProductID {
				uberCab.Arriving = cabTimes.Estimate / 60 // Response time is in seconds
				break
			}
		}
		uberCabs = append(uberCabs, uberCab)
	}
	return uberCabs, nil
}

func uberRequest(url string) ([]byte, error) {
	reqPrices, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	reqPrices.Header.Set("Content-Type", "application/json")
	reqPrices.Header.Set("Authorization", os.Getenv("UBER_AUTH_TOKEN"))
	reqPrices.Header.Set("Accept-Language", "en_US")

	client := http.Client{Timeout: time.Second * 60}

	respPrices, err := client.Do(reqPrices)
	if err != nil {
		return nil, err
	}
	defer func() { _ = respPrices.Body.Close() }()

	respBytes, err := ioutil.ReadAll(respPrices.Body)
	return respBytes, err
}
