package scraper

import (
	"context"
	"errors"
	docclient "github.com/david-callender/FoodFinder/dineocclient"
	"math/rand"
	"time"
	pgx "github.com/jackc/pgx/v5"
	"strings"
)

// Global variables

// I'd like this to be a constant but constant arrays don't exist in go.
var mealtimeIndexer [4]string = [4]string{"breakfast", "lunch", "dinner", "everyday"}

const PAST_SCRAPE_DAYS int = 7
const SCRAPE_PERIOD int = 14
const SLEEP_MIN_SECS int = 5
const SLEEP_MAX_DIFF int = 10
const TIME_DAY time.Duration = 24 * time.Hour

// NON-EXPORTED FUNCTIONS

// scrapeMenuToDatabase(conn (*pgx.Conn), locationId, periodName (strings), date (time.Time):
// Takes a database connection, a dineoncampus location ID, a meal period name,
// and a date. It removes all old menu data for the menu corresponding to these
// parameters, and fills in new menu data if it exists.
func scrapeMenuToDatabase(conn *pgx.Conn, locationId, periodName string, date time.Time) error {
	// Fetch the menu data for insertion into the DB
	menu, err := docclient.GetMenuById(locationId, periodName, date)
	if err != nil {
		return err
	}

	// Convert periodName to a 16-bit integer identifier as used in our db.
	// Int16 to avoid type issues because mealtime is a smallint in our db.
	var mealtimeNum int16 = 255 // 255 stands in for our uninitialized value.
	for i, timeName := range mealtimeIndexer {
		if strings.ToLower(periodName) == timeName {
			mealtimeNum = int16(i)
		}
	}
	if mealtimeNum == 255 {
		return errors.New("scrapeMenuToDatabase: Invalid period name.")
	}
	
	dateFormatted := date.Format("2006-01-02")

	// Scraping a location-period-day menu into the db updates the db by
	// first deleting all of the old menu data (in case it has changed).
	// We use a transaction to prevent failed updates from leaving a day's
	// menu completely empty or in an inconsistent state.
	transaction, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	// This can be safely done since transaction.Rollback returns an error
	// once the transaction has been closed.
	defer transaction.Rollback(context.Background())
	
	_, err = transaction.Exec(
		context.Background(),
		`DELETE FROM "DocCache" WHERE day=$1 AND location=$2 AND mealtime=$3;`,
		dateFormatted, locationId, mealtimeNum,
	)
	if err != nil {
		return err
	}

	// Only perform the insert operation if there's stuff to insert.
	if len(menu.Options) != 0 {
		_, err = transaction.CopyFrom(
			context.Background(),
			pgx.Identifier{"DocCache"},
			[]string{"day", "location", "mealtime", "meal", "mealid"},
			pgx.CopyFromSlice(len(menu), func(i int) ([]any, error) {
				return []any{
					dateFormatted,
					locationId,
					mealtimeNum,
					menu.Options[i].Name,
					menu.Options[i].Id,
				}, nil
			}),
		)
		if err != nil {
			return err
		}
	}

	if err := transaction.Commit(context.Background()); err != nil {
		return err
	}

	// If all goes well, we simply return no error.
	return nil
}

// EXPORTED FUNCTIONS FOR USE BY IMPORTING CODE

// ScrapeMenusToDatabase(conn (*pgx.Conn), siteId (string)): Takes a database
// connection and a dineoncampus site ID. It will scrape all available menus
// into the database.
func ScrapeMenusToDatabase(conn *pgx.Conn, siteId string) error {
	// Fetch all food locations at a given site
	foodBuildings, err := docclient.GetFoodBuildings(siteId)
	if err != nil {
		return err
	}

	// We need a list of times to scrape, for now this is a week before and
	// a week after the current date. Set global constants to change this.
	var dates []time.Time
	currentDate := time.Now().Truncate(TIME_DAY)
	for i := 0; i <= SCRAPE_PERIOD; i++ {
		timeOffset := (time.Duration(i) * TIME_DAY) - (time.Duration(PAST_SCRAPE_DAYS) * TIME_DAY)
		dates = append(dates, currentDate.Add(timeOffset))
	}

	// Delete all menus that are older than 1 week to prevent endlessly growing
	// the db.
	_, err = conn.Exec(
		context.Background(),
		`DELETE FROM "DocCache" WHERE day < $1`,
		dates[0].Format("2006-01-02"),
	)
	if err != nil {
		return err
	}

	for _, building := range foodBuildings {
		for _, location := range building.Locations {
			for _, periodName := range mealtimeIndexer {
				for _, date := range dates {
					err = scrapeMenuToDatabase(conn, location.Id, periodName, date)
					if err != nil {
						return err
					}
					// We sleep between SLEEP_MIN_SECS and
					// SLEEP_MAX_DIFF + SLEEP_MIN_SECS seconds
					// to prevent rate limiting or overloading
					// our database when scraping.
					time.Sleep(time.Duration(rand.Intn(SLEEP_MAX_DIFF) + SLEEP_MIN_SECS) * time.Second)
				}
			}
		}
	}
	
	return nil
}
