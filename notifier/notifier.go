package notifier

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

func main() {
	notifyDate := time.Now().Truncate(time.Duration(24) * time.Hour)

	connString := os.Getenv("DATABASE_URL")
	conn := pgx.Connect(context.Background(), connString)

	dateFormatted := notifyDate.Format("2006-01-02")
	usersToNotify := conn.Query(
		context.Background(),
		`SELECT user, meal 
			FROM "Preferences"
			JOIN "DocCache" 
			ON "Preferences.preference" = "DocCache.meal" 
			WHERE day=$1
			ORDER BY user;`,
		dateFormatted,
	)

	var emailTable map[int]string
	userEmails := conn.Query(
		context.Background(),
		`SELECT id, email FROM "Users" JOIN "Preferences" ON "Users.id" = "Preferences.user";`,
	)
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

}
