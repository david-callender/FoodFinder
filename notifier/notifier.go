package notifier

import (
	"context"
	"os"
	"time"

	docclient "github.com/david-callender/FoodFinder/dineocclient"
	"github.com/jackc/pgx/v5"
	"gopkg.in/gomail.v2"
)

// Global Constant Storage
// Generic Global Values
const NOTIFIER_EMAIL string = "<EMAIL>@<DOMAIN>.<EXTENSION>"
const NOTIFICATION_SUBJECT string = "GopherGrub Notification"
const TIME_DAY time.Duration = 24 * time.Hour

func main() {
	notifyDate := time.Now().Truncate(time.Duration(24) * time.Hour)

	connString := os.Getenv("DATABASE_URL")
	conn := pgx.Connect(context.Background(), connString)

}

func notifyUsers(conn *pgx.Conn, date time.Time) error {
	date = date.Truncate(TIME_DAY)

	dateFormatted := date.Format("2006-01-02")
	usersToNotify, err := conn.Query(
		context.Background(),
		`SELECT user, meal 
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

	var notificationTable map[int][]string
	for usersToNotify.Next() {
		// Again we can't short declare since Scan takes a reference.
		var userId int
		var meal string
		usersToNotify.Scan(&userId, &meal)
		notificationTable[userId] = append(notificationTable[userId], meal)
	}

	var messages []gomail.Message
	i := 0
	for userId, meals := range notificationTable {
		message := gomail.Message{}
		message.SetHeader("From", NOTIFIER_EMAIL)
		message.SetHeader("To", emailTable[userId])
		message.SetHeader("Subject", NOTIFICATION_SUBJECT)
		messageBody := "Some of your favorite foods are available today!"
		for j, meal := range meals {

		}
	}

}
