package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	docclient "github.com/david-callender/FoodFinder/notifier/dineocclient"
	"github.com/jackc/pgx/v5"
	"github.com/wneessen/go-mail"
)

// Global Constant Storage
// Generic Global Values
var mealtimeIndexer = [4]string{"breakfast", "lunch", "dinner", "every day"}

const EMAIL_HOST = "<EMAIL_SERVER_HOSTNAME>"
const NOTIFICATION_SUBJECT = "GopherGrub Notification"
const TIME_DAY = 24 * time.Hour
const UMN_SITE_ID = "61d7515eb63f1e0e970debbei"

// Types
type mealNotification struct {
	user     int
	meal     string
	location string
	mealTime int16
}

// Functions

// The main function. Reads the database connection string from the environment,
// and the date for which to send notifications from the commandline. Calls
// the notifyUsers function with the extracted values.
func main() {
	connString := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalln(err)
	}

	notifyTime, err := time.Parse("2006-01-02", os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	err = notifyUsers(conn, notifyTime)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}

// generateMessages(notificationTable, emailTable): Takes a mapping from integer
// user ids to a slice of meals the user needs to be notified of. Also takes a
// mapping from integer user ids to their emails. Prepares a slice of *mail.Msg
// structs to be sent out to users. Returns a nil slice and non-nil error on
// failure.
func generateMessages(notificationTable map[int][]mealNotification, emailTable map[int]string) ([]*mail.Msg, error) {
	locationIdToName, err := getLocations(UMN_SITE_ID)
	if err != nil {
		return nil, err
	}

	var messages = make([]*mail.Msg, 0)
	for userId, notifs := range notificationTable {
		message := mail.NewMsg()
		errs := errors.Join(
			message.From(os.Getenv("NOTIFIER_EMAIL")),
			message.ToFromString(emailTable[userId]),
		)
		if errs != nil {
			return nil, errs
		}
		message.Subject(NOTIFICATION_SUBJECT)
		message.SetBulk()
		messageBody := "Some of your favorite foods are available today!\n\n"
		for _, notif := range notifs {
			messageBody += fmt.Sprintf(
				"- %s at %s during %s time.\n",
				notif.meal,
				locationIdToName[notif.location],
				mealtimeIndexer[notif.mealTime],
			)
		}
		message.SetBodyString("text/plain", messageBody)
		messages = append(messages, message)
	}

	return messages, nil
}

// getLocations(siteId): Takes a 24-character Dineoncampus site id and returns
// a mapping from the location id string to the location name. Returns an
// empty map and non-nil error on failure.
func getLocations(siteId string) (map[string]string, error) {
	var locationIdToName = make(map[string]string)
	buildings, err := docclient.GetFoodBuildings(siteId)
	if err != nil {
		return locationIdToName, err
	}
	for _, building := range buildings {
		for _, location := range building.Locations {
			locationIdToName[location.Id] = location.Name
		}
	}
	return locationIdToName, err
}

// notifyUsers(conn, date): takes a pgx database connection and a date (as time.Time)
// for which to notify users. Obtains a list of matches between user preferences
// and meals in the cache, and a mapping between integer user ids and their emails.
// Then it will call generateMessages() to obtain a list of messages and send them
// with sendMessages(). Returns non-nil error on failure.
func notifyUsers(conn *pgx.Conn, date time.Time) error {
	date = date.Truncate(TIME_DAY)

	dateFormatted := date.Format("2006-01-02")
	usersToNotify, err := conn.Query(
		context.Background(),
		`SELECT user, meal, location, mealtime 
			FROM "Preferences"
			JOIN "DocCache" 
			ON "Preferences.preference" = "DocCache.meal" 
			WHERE day=$1
			ORDER BY user;`,
		dateFormatted,
	)
	if err != nil {
		return err
	}

	var emailTable = make(map[int]string)
	userEmails, err := conn.Query(
		context.Background(),
		`SELECT id, email FROM "Users" JOIN "Preferences" ON "Users.id" = "Preferences.user";`,
	)
	if err != nil {
		return err
	}
	for userEmails.Next() {
		// These cannot be short-declared since Scan takes a reference.
		var userId int
		var email string
		err = userEmails.Scan(&userId, &email)
		if err != nil {
			return err
		}
		emailTable[userId] = email
	}

	var notificationTable = make(map[int][]mealNotification)
	for usersToNotify.Next() {
		// Again we can't short declare since Scan takes a reference.
		var notif mealNotification
		err = usersToNotify.Scan(&notif.user, &notif.meal, &notif.location, &notif.mealTime)
		if err != nil {
			return err
		}
		notificationTable[notif.user] = append(notificationTable[notif.user], notif)
	}

	messages, err := generateMessages(notificationTable, emailTable)
	if err != nil {
		return err
	}

	if err = sendMessages(messages); err != nil {
		return err
	}

	return nil
}

// sendMessages(messages): takes a slice of references to mail.Msg emails, obtains
// its email address and password from the environment, and uses the credentials
// to send all of the messages. Returns non-nil error on failure.
func sendMessages(messages []*mail.Msg) error {
	mailer, err := mail.NewClient(
		EMAIL_HOST,
		mail.WithUsername(os.Getenv("NOTIFIER_EMAIL")),
		mail.WithPassword(os.Getenv("NOTIFIER_PASSWORD")),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return err
	}

	if err = mailer.DialAndSend(messages...); err != nil {
		return err
	}

	return nil
}
