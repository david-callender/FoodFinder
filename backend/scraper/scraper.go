package "scraper"

import (
	"errors",
	docclient "github.com/david-callender/foodfinder/dineocclient",
	"fmt",
	"github.com/jackc/pgx",
	"sql",
)

// NON-EXPORTED FUNCTIONS

func scrapeMenuToDatabase(locationId, periodName string, date time.Time) error {
	menu, err := docclient.GetMenuById(locationId, periodName, date)
	if err != nil {
		return err
	}

}

// EXPORTED FUNCTIONS FOR USE BY IMPORTING CODE

func ScrapeMenusToDatabase(siteId string) {

}
