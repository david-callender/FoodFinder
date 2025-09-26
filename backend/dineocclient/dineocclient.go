package dineocclient

import (
	jsonv2 "encoding/json/v2";
//	"fmt";
	"io";
	"net/http";
//	"os";
	"time"
)

// Global config variables that we ought to move out to a config file
const dineocaddress = "https://apiv4.dineoncampus.com/"
const useragent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0"

// EXPORTED TYPES: MAY BE USED BY IMPORTING MODULES

type FoodBuilding struct {
	Building string `json:"buildingName"`
	Locations []Restaurant `json:"locations"`
}

type Meal struct {
	Name, Description string
}

type Menu struct {
	Period string
	Options []Meal
}

type Restaurant struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// NON-EXPORTED TYPES: MAY NOT BE USED BY IMPORTING MODULES

type PeriodIdSpec struct {
	Breakfast, Dinner, Lunch, Everyday string
}

type period struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// EXPORTED FUNCTIONS: MAY BE USED BY IMPORTING MODULES

// GetFoodBuildings(site string): Takes a dineoncampus site id and returns a
// slice of FoodLocation structs representing all food locations found for
// the site.
func GetFoodBuildings(site string) ([]FoodBuilding, error) {
	apifunc := dineocaddress + "sites/" + site + "/locations-public?for_menus=true"
	jsonData, err := makeDineocApiCall(apifunc)
	if err != nil {
		return nil, err
	}
	
	var buildings struct {
		Buildings []FoodBuilding `json:"buildings"`
	}
	if err := jsonv2.Unmarshal([]byte(jsonData), &buildings); err != nil {
		return nil, err
	}

	return buildings.Buildings, nil
}

// NON-EXPORTED FUNCTIONS: MAY NOT BE USED BY IMPORTING MODULES

// getPeriodIds(location, date): takes a dining hall location ID and returns
// a periodIdSpec containing a set of meal period IDs for breakfast, lunch,
// dinner, and everyday menus. These IDs may only be used once each.
func GetPeriodIds(location string, date time.Time) (PeriodIdSpec, error) {
	// This has to be defined up here in case of an error so we have a 0 value
	var periodIds PeriodIdSpec

	dateFormatted := date.Format("2006-01-02")
	apifunc := dineocaddress + "locations/" + location + "/periods/?date=" + dateFormatted
	jsonData, err := makeDineocApiCall(apifunc)
	if err != nil {
		return periodIds, err
	}

	var periods struct {
		Periods []period `json:"periods"`
	}
	if err := jsonv2.Unmarshal([]byte(jsonData), &periods); err != nil {
		return periodIds, err
	}
	
	periodIds = PeriodIdSpec{
		Breakfast: periods.Periods[0].Id,
		Lunch: periods.Periods[1].Id,
		Dinner: periods.Periods[2].Id,
		Everyday: periods.Periods[3].Id,
	}

	return periodIds, nil
}

// makeDineocApiCall(apiurl): make a GET request to apiurl, and return the
// body of the response.
func makeDineocApiCall(apiurl string) ([]byte, error) {
	req, err := newDineocApiRequest(apiurl, "GET")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// newDineocApiRequest(apiurl, method): return a http request to the given apiurl
// with properly populated headers for making a request to the dineoc api.
func newDineocApiRequest(apiurl string, method string) (*http.Request, error) {
	req, err := http.NewRequest("GET", apiurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", useragent)
	req.Header.Add("accept", "application/json")
	// req is already a pointer because http.newRequest returns a pointer
	return req, nil
}
