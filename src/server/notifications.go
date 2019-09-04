package server

import (
	"bytes"
	"database/sql"
	"text/template"
	"time"

	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func createNotificationTmpl(db *sql.DB, department string) string {
	notifications, _ := lmdatabase.QueryNotifications(db, department)
	funcMap := template.FuncMap{
		"priorityMap": func(prio int) string {
			priorityMap := map[int]string{
				1: "is-danger",
				2: "is-warning",
				3: "is-info",
			}
			return priorityMap[prio]
		},
		"now": func(now int64) string {
			return time.Unix(now, 0).Format("15:04:05")
		},
	}
	data := map[string]interface{}{
		"Notifications": notifications,
	}
	var notificationsBuffer bytes.Buffer
	notificationHTML := template.Must(template.New("notifications.html").Funcs(funcMap).ParseFiles("templates/notifications.html"))
	notificationHTML.Execute(&notificationsBuffer, data)
	return notificationsBuffer.String()
}
