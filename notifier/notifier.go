package notifier

import (
	"context"
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
var mealtimeIndexer [4]string = [4]string{"breakfast", "lunch", "dinner", "every day"}

const EMAIL_HOST string = "<EMAIL_SERVER_HOSTNAME>"
const NOTIFIER_EMAIL string = "<EMAIL>@<DOMAIN>.<EXTENSION>"
const NOTIFICATION_SUBJECT string = "GopherGrub Notification"
const TIME_DAY time.Duration = 24 * time.Hour
const UMN_SITE_ID string = "61d7515eb63f1e0e970debbei"

// Types
type mealNotification struct {
	user     int
	meal     string
	location string
	mealTime int16
}

func main() {
	connString := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}

	notifyTime, err := time.Parse("2006-01-02", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	notifyUsers(conn, notifyTime)

	return
}

func generateMessages(notificationTable map[int][]mealNotification, emailTable map[int]string) ([]*mail.Msg, error) {
	locationIdToName, err := getLocations(UMN_SITE_ID)
	if err != nil {
		return nil, err
	}

	var messages []*mail.Msg
	for userId, notifs := range notificationTable {
		message := mail.NewMsg()
		message.From(NOTIFIER_EMAIL)
		message.ToFromString(emailTable[userId])
		message.Subject(NOTIFICATION_SUBJECT)
		message.SetBulk()
		messageBody := fmt.Sprintf("Some of your favorite foods are available today!\n\n")
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

func getLocations(siteId string) (map[string]string, error) {
	var locationIdToName map[string]string
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

	var emailTable map[int]string
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
		userEmails.Scan(&userId, &email)
		emailTable[userId] = email
	}

	var notificationTable map[int][]mealNotification
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

func sendMessages(messages []*mail.Msg) error {
	mailer, err := mail.NewClient(
		EMAIL_HOST,
		mail.WithUsername(NOTIFIER_EMAIL),
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
