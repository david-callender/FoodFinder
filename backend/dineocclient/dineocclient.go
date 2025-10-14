package dineocclient

import (
	jsonv2 "encoding/json/v2";
	"fmt";
	"io";
	"net/http";
	"math/rand";
	"strings";
	"time"
)

// Global config variables that we ought to move out to a config file
const dineocaddress = "https://apiv4.dineoncampus.com/"

// EXPORTED TYPES: MAY BE USED BY IMPORTING MODULES

type FoodBuilding struct {
	Name string `json:"buildingName"`
	Locations []Restaurant `json:"locations"`
}

type Meal struct {
	Id string `json:"id"`
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
// Periods are themselves returned as json objects, which requires a specific
// struct, even though our periodIdSpec simply has four fields named after the
// period ID in question.
type period struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// BEGIN: Group of structs used for parsing the menu data

// Dineoncampus separates meals by category, which does not seem to have any
// functional purpose. Still, we need go structs that match the shape of the
// json data returned by dineoncampus.
type category struct {
	Items []Meal `json:"items"`
}

// Same deal here, a menu is returned under a period, with the foods being in
// categories under the menu. Since our data isn't shaped like this, we need
// another struct to make sure it can be properly unmarshaled.
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
	
	// This struct is not anonymous as jsonv2.Unmarshal takes a reference to it
	// We need the data contained within later, thus it must be a variable.
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
	var locationId string = ""

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

// GetMenuByName(buildingName, locationName, periodName, site (strings), date (time)): 
// Takes the name of a building, the name of a location within the building,
// a named meal period ("breakfast", "lunch", "dinner", or "everyday"), a
// dineoncampus site ID, and a time.Time representing the date for the menu requested. 
// Returns a Menu populated with the options from dineoncampus. Names are
// case-insensitive.
func GetMenuByName(buildingName, locationName, periodName, siteId string, date time.Time) (Menu, error) {
	locationId, err := GetLocationIdByName(buildingName, locationName, siteId)
	if err != nil {
		return Menu{}, err
	}

	return GetMenuById(locationId, periodName, date)
}

// GetMenuById(locationId, periodName (strings), date (time.Time)): Takes the
// ID of a food location (NOT A BUILDING ID), a named meal period, and a time.Time
// representing a calendar date. Returns a Menu populated with the options from
// dineoncampus.
func GetMenuById(locationId, periodName string, date time.Time) (Menu, error) {
	var menu Menu = Menu{Locations: []Restaurant{}}
	var periodId string = ""

	dateFormatted := date.Format("2006-01-02")
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
		return menu, nil
	}

	apifunc := dineocaddress +
		"locations/" + locationId +
		"/menu?date=" + dateFormatted +
		"&period=" + periodId
	jsonData, err := makeDineocApiCall(apifunc)
	if err != nil {
		return menu, err
	}

	// This struct is not anonymous as jsonv2.Unmarshal takes a reference to it
	// We need the data contained within later, thus it must be a variable.
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

	// Since we get an array of categories which contain an array of meals,
	// two loops are required to flatten the two lists into one big list of
	// meals that are not separated by category.
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
	var periodIds periodIdSpec = periodIdSpec{}

	dateFormatted := date.Format("2006-01-02")
	apifunc := dineocaddress +
		"locations/" + locationId +
		"/periods/?date=" + dateFormatted
	jsonData, err := makeDineocApiCall(apifunc)
	if err != nil {
		return periodIds, err
	}

	// This struct is not anonymous as jsonv2.Unmarshal takes a reference to it
	// We need the data contained within later, thus it must be a variable.
	var periods struct {
		Periods []period `json:"periods"`
	}
	if err := jsonv2.Unmarshal([]byte(jsonData), &periods); err != nil {
		return periodIds, err
	}
	
	// We can't assume that there will always be all four periods, because
	// sometimes one or more periods are not returned. When this happens,
	// it also moves all of the other periods indexes. Thus we have to loop
	// through and check the name associated to each ID to determine which
	// field it belongs to.
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
	// We have to add some minimum of headers to ensure that we get the right
	// data, and also to make sure that the client does not get blocked by
	// cloudflare (typically due to a bad useragent).
	req.Header.Add("user-agent", useragents[rand.Intn(len(useragents))])
	req.Header.Add("accept", "application/json")
	// req is already a pointer because http.newRequest returns a pointer
	return req, nil
}
