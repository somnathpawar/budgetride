package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
	"time"
	"os"
)

func LyftCabs(startLn, startLg, endLn, endLg string) (CabList, error) {
    type cabPrices struct {
        Currency        string  `json:"currency"`
        RideType        string  `json:"ride_type"`
        DisplayName     string  `json:"display_name"`
        PrimePercentage string  `json:"primetime_percentage"`
        Duration        int     `json:"estimated_duration_seconds"`
        Miles           float64 `json:"estimated_distance_miles"`
        MinCost         int     `json:"estimated_cost_cents_min"`
        MaxCost         int     `json:"estimated_cost_cents_max"`
        CanRequestRide  bool    `json:"can_request_ride"`
    }
    type estimatePrices struct {
        Prices []cabPrices `json:"cost_estimates"`
    }
    priceBody := estimatePrices{}
    lyftPrices, err := lyftRequest("https://api.lyft.com/v1/cost?start_lat=" + startLn + "&start_lng=" + startLg + "&end_lat=" + endLn + "&end_lng=" + endLg)
    err = json.Unmarshal(lyftPrices, &priceBody)
    if err != nil {
        return nil, err
    }

    type cabTimes struct {
        RideType    string `json:"ride_type"`
        Estimate    int    `json:"eta_seconds"`
        DisplayName string `json:"display_name"`
    }
    type estimateTimes struct {
        Times []cabTimes `json:"eta_estimates"`
    }
    estimateBody := estimateTimes{}
    lyftTimes, err := lyftRequest("https://api.lyft.com/v1/eta?lat=" + startLn + "&lng=" + startLg)
    err = json.Unmarshal(lyftTimes, &estimateBody)
    if err != nil {
        return nil, err
    }

    lyftCabs := CabList{}
    for _, cab := range priceBody.Prices {
        lyftCab := CabResponse{}
        lyftCab.Company = "Lyft"
        lyftCab.Cab = cab.DisplayName
        lyftCab.Estimate = float64(cab.MinCost / 100)
        for _, cabTimes := range estimateBody.Times {
            if cabTimes.RideType == cab.RideType {
                lyftCab.Arriving = cabTimes.Estimate / 60 // Response time is in seconds
                break
            }
        }

        lyftCabs = append(lyftCabs, lyftCab)
    }
    return lyftCabs, nil
}

func lyftRequest(url string) ([]byte, error) {
    reqPrices, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Set headers
    reqPrices.Header.Set("Content-Type", "application/json")
    reqPrices.Header.Set("Authorization", os.Getenv("LYFT_AUTH_TOKEN"))
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