package dineocclient

import (
	jsonv2 "encoding/json/v2";
	"errors"
	"fmt";
	"io";
	"net/http";
//	"os";
	"strings";
	"time"
)

// Global config variables that we ought to move out to a config file
const dineocaddress = "https://apiv4.dineoncampus.com/"
const useragent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0"

// EXPORTED TYPES: MAY BE USED BY IMPORTING MODULES

type FoodBuilding struct {
	Name string `json:"buildingName"`
	Locations []Restaurant `json:"locations"`
}

type Meal struct {
	Name string `json:"name"`
	Description string `json:"desc"`
}

type Menu struct {
	Date time.Time
	PeriodName string
	Options []Meal
}

type Restaurant struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// NON-EXPORTED TYPES: MAY NOT BE USED BY IMPORTING MODULES

// Internal struct used for passing one-time use period IDs
type periodIdSpec struct {
	Breakfast, Dinner, Lunch, Everyday string
}

// NON-EXPORTED STRUCTS USED ONLY FOR PARSING JSON

// Period struct used for parsing period IDs into a periodIdSpec
type period struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// BEGIN: Group of structs used for parsing the menu data
type category struct {
	Items []Meal `json:"items"`
}

type menuPeriod struct {
	Name string `json:"name"`
	Categories []category `json:"categories"`
}
// END: Group of structs used for parsing the menu data

// EXPORTED FUNCTIONS: MAY BE USED BY IMPORTING MODULES

// GetFoodBuildings(site string): Takes a dineoncampus site id and returns a
// slice of FoodLocation structs representing all food locations found for
// the site.
func GetFoodBuildings(siteId string) ([]FoodBuilding, error) {
	apifunc := dineocaddress +
		"sites/" + siteId +
		"/locations-public?for_menus=true"
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

// GetLocationIdByName(buildingName, locationName): Takes the name of a building,
// the name of a location, and a dineoncampus site ID, and returns the location
// ID corresponding to that location. Names are case-insensitive
func GetLocationIdByName(buildingName, locationName, siteId string) (string, error) {
	var locationId string

	buildingName = strings.ToLower(buildingName)
	locationName = strings.ToLower(locationName)

	buildings, err := GetFoodBuildings(siteId)
	if err != nil {
		return "", err
	}

	for _, building := range buildings {
		building.Name = strings.ToLower(building.Name)
		if building.Name != buildingName {
			continue // Skip the rest if we aren't in the right building
		}
		for _, location := range building.Locations {
			if strings.ToLower(location.Name) == locationName {
				locationId = location.Id
			}
		}
	}

	if locationId == "" {
		err := fmt.Errorf("GetLocationIdByName(): Could not find location: %v", locationName)
		return "", err
	}

	return locationId, nil
}

// GetMenu(buildingName, locationName, periodName, site (strings), date (time)): 
// Takes the name of a building, the name of a location within the building,
// a named meal period ("breakfast", "lunch", "dinner", or "everyday"), a
// dineoncampus site ID, and a time.Time representing the current date. 
// Returns a Menu populated with the options from dineoncampus. Names are
// case-insensitive.
func GetMenu(buildingName, locationName, periodName, siteId string, date time.Time) (Menu, error) {
	var menu Menu
	var periodId string

	dateFormatted := date.Format("2006-01-02")
	locationId, err := GetLocationIdByName(buildingName, locationName, siteId)
	if err != nil {
		return menu, err
	}
	periodIds, err := getPeriodIds(locationId, date)
	if err != nil {
		return menu, err
	}

	switch strings.ToLower(periodName) {
	case "breakfast":
		periodId = periodIds.Breakfast
	case "lunch":
		periodId = periodIds.Lunch
	case "dinner":
		periodId = periodIds.Dinner
	case "everyday":
		periodId = periodIds.Everyday
	default:
		err := fmt.Errorf("GetMenu(): Invalid period name: %v.", periodName)
		return menu, err
	}

	if periodId == "" {
		err := errors.New("GetMenu(): Failed to get period ID")
		return menu, err
	}

	apifunc := dineocaddress +
		"locations/" + locationId +
		"/menu?date=" + dateFormatted +
		"&period=" + periodId
	jsonData, err := makeDineocApiCall(apifunc)
	if err != nil {
		return menu, err
	}

	var rawMenu struct {
		Period menuPeriod `json:"period"`
	}
	if err := jsonv2.Unmarshal([]byte(jsonData), &rawMenu); err != nil {
		return menu, err
	}

	menu.Date = date
	menu.Options = make([]Meal, 0, 25)
	menu.PeriodName = rawMenu.Period.Name
	menuCategories := rawMenu.Period.Categories

	for _, category := range menuCategories {
		for _, 	mealOption := range category.Items {
			menu.Options = append(menu.Options, mealOption)
		}
	}

	return menu, nil
}

// NON-EXPORTED FUNCTIONS: MAY NOT BE USED BY IMPORTING MODULES

// getPeriodIds(locationId, date): takes a dining hall location ID and returns
// a periodIdSpec containing a set of meal period IDs for breakfast, lunch,
// dinner, and everyday menus. These IDs may only be used once each.
func getPeriodIds(locationId string, date time.Time) (periodIdSpec, error) {
	// This has to be defined up here in case of an error so we have a 0 value
	var periodIds periodIdSpec

	dateFormatted := date.Format("2006-01-02")
	apifunc := dineocaddress +
		"locations/" + locationId +
		"/periods/?date=" + dateFormatted
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
	
	for _, period := range periods.Periods {
		switch strings.ToLower(period.Name) {
		case "breakfast":
			periodIds.Breakfast = period.Id
		case "lunch":
			periodIds.Lunch = period.Id
		case "dinner":
			periodIds.Dinner = period.Id
		case "every day":
			periodIds.Everyday = period.Id
		}
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
