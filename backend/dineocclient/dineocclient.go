package dineocclient

import (
	jsonv2 "encoding/json/v2";
	"fmt";
	"io";
	"net/http";
	"os"
)

// Global config variables that we ought to move out to a config file
const dineocaddress = "https://apiv4.dineoncampus.com/"
const useragent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0"
const umnsiteid = "61d7515eb63f1e0e970debbe"

// EXPORTED TYPES: MAY BE USED BY IMPORTING MODULES

type FoodLocation struct {
	Id, Name string
}

type Meal struct {
	Name, Description string
}

type Menu struct {
	Period string
	Options []Meal
}

// NON-EXPORTED TYPES: MAY NOT BE USED BY IMPORTING MODULES



// EXPORTED FUNCTIONS: MAY BE USED BY IMPORTING MODULES

// GetFoodLocation(site string): Takes a dineoncampus site id and returns a
// slice of FoodLocation structs representing all food locations found for
// the site.
func GetFoodLocations(site string) ([]FoodLocation, error) {
	apifunc := dineocaddress + "sites/" + site + "/locations-public?for-menus=true"
	req, err := http.NewRequest("GET", apifunc, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", useragent)
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

// NON-EXPORTED FUNCTIONS: MAY NOT BE USED BY IMPORTING MODULES


