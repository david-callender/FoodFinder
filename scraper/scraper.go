package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	docclient "github.com/david-callender/FoodFinder/utils/dineocclient"

	"github.com/jackc/pgx/v5"
)

// Global variables

// I'd like this to be a constant but constant arrays don't exist in go.
var mealtimeIndexer [4]string = [4]string{"breakfast", "lunch", "dinner", "everyday"}
var hallsToScrape [6]string = [6]string{
	"Comstock Dining Hall",
	"17th Ave. Dining Hall",
	"Pioneer Dining Hall",
	"Sanford Dining Hall",
	"Middlebrook Dining Hall",
	"Bailey Dining Hall",
}

var errNoConnString error = errors.New("DATABASE_URL is not set, cannot connect to database")
var errNoScrapeBackArg error = errors.New("-back requires an argument")
var errNoScrapeFwArg error = errors.New("-forward requires an argument")

const DEFAULT_PAST_SCRAPE int = 7
const DEFAULT_FUTURE_SCRAPE int = 14
const SLEEP_MIN_SECS int = 5
const SLEEP_MAX_DIFF int = 10
const TIME_DAY time.Duration = 24 * time.Hour
const UMN_SITE_ID string = "61d7515eb63f1e0e970debbe"

// main(): calls runScraper() and logs any errors received, logging success otherwise.
func main() {
	err := runScraper()
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		log.Println("Successfully scraped menus to database")
	}
}

// NON-EXPORTED FUNCTIONS
func getLocations(siteId string) ([]docclient.Restaurant, error) {
	foodBuildings, err := docclient.GetFoodBuildings(siteId)
	if err != nil {
		return nil, err
	}

	var locations []docclient.Restaurant
	for _, building := range foodBuildings {
		for _, location := range building.Locations {
			for _, name := range hallsToScrape {
				if location.Name == name {
					locations = append(locations, location)
				}
			}
		}
	}

	return locations, nil
}

// scrapeMenuToDatabase(conn (*pgx.Conn), locationId, periodName (strings), date (time.Time):
// Takes a database connection, a dineoncampus location ID, a meal period name,
// and a date. It removes all old menu data for the menu corresponding to these
// parameters, and fills in new menu data if it exists.
func scrapeMenuToDatabase(conn *pgx.Conn, locationId, periodName string, date time.Time) error {
	// Predefining the nil error to ensure err exists
	var err error

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
		return errors.New("scrapeMenuToDatabase: Invalid period name")
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
	defer func() {
		err = transaction.Rollback(context.Background())
		if err == pgx.ErrTxClosed {
			err = nil
		}
	}()

	_, err = transaction.Exec(
		context.Background(),
		`DELETE FROM "DocCache" WHERE day=$1 AND location=$2 AND mealtime=$3`,
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
			pgx.CopyFromSlice(len(menu.Options), func(i int) ([]any, error) {
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
	return err
}

// runScraper(): Reads the database URL from the environment, and optionally
// a number of dates to scrape backward and forward from the current date. It
// then calls ScrapeMenusToDatabase() with the collected values.
func runScraper() error {
	var err error
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return errNoConnString
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close(context.Background())
	}()

	dates := make([]time.Time, 0)
	currentDate := time.Now().Truncate(time.Hour)
	scrapeBackDays := DEFAULT_PAST_SCRAPE
	scrapeForwardDays := DEFAULT_FUTURE_SCRAPE
	for i, arg := range os.Args {
		switch arg {
		case "-back":
			if len(os.Args) < (i + 2) {
				return errNoScrapeBackArg
			}
			intArg, err := strconv.ParseInt(os.Args[i+1], 10, 0)
			scrapeBackDays = int(intArg)
			if err != nil {
				return err
			}
		case "-forward":
			if len(os.Args) < (i + 2) {
				return errNoScrapeFwArg
			}
			intArg, err := strconv.ParseInt(os.Args[i+1], 10, 0)
			if err != nil {
				return err
			}
			scrapeForwardDays = int(intArg)
		}
	}
	scrapePeriod := scrapeBackDays + scrapeForwardDays
	for i := 0; i <= scrapePeriod; i++ {
		timeOffset := (time.Duration(i) * TIME_DAY) - (time.Duration(scrapeBackDays) * TIME_DAY)
		dates = append(dates, currentDate.Add(timeOffset))
	}

	log.Printf("Scraping from %d day(s) ago to %d day(s) from now\n", scrapeBackDays, scrapeForwardDays)
	err = ScrapeMenusToDatabase(conn, dates, UMN_SITE_ID)
	return err
}

// EXPORTED FUNCTIONS FOR USE BY IMPORTING CODE

// ScrapeMenusToDatabase(conn (*pgx.Conn), dates ([]time.Time, siteId (string)):
// Takes a database connection, a dineoncampus site ID, and a list of dates to scrape.
// Will scrape all available menus for the given days into the database.
func ScrapeMenusToDatabase(conn *pgx.Conn, dates []time.Time, siteId string) error {
	// Fetch all food locations at a given site
	locations, err := getLocations(siteId)
	if err != nil {
		return err
	}

	// Delete all menus that are older than the oldest given date to prevent
	// endlessly growing the db.
	_, err = conn.Exec(
		context.Background(),
		`DELETE FROM "DocCache" WHERE day < $1`,
		dates[0].Format("2006-01-02"),
	)
	if err != nil {
		return err
	}

	for _, location := range locations {
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
				time.Sleep(time.Duration(rand.Intn(SLEEP_MAX_DIFF)+SLEEP_MIN_SECS) * time.Second)
			}
		}
	}
	return nil
}
