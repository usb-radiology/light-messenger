package server

import (
	"bytes"
	"database/sql"
	"log"
	"text/template"
	"time"

	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func getCardHTML(db *sql.DB, modality string, department string) string {

	notification, _ := lmdatabase.NotificationGetByDepartmentAndModality(db, department, modality)
	// AOD Card Template
	data := map[string]interface{}{
		"Modality":       notification.Modality,
		"Department":     notification.DepartmentID,
		"PriorityNumber": notification.Priority,
		"PriorityName":   priorityMap[notification.Priority],
		"CreatedAt":      time.Unix(notification.CreatedAt, 0).Format("15:04:05"),
	}

	var aodBuffer bytes.Buffer
	templates[templateCardID].Execute(&aodBuffer, data)
	return aodBuffer.String()
}

func getNotificationHTML(db *sql.DB, department string) string {
	notifications, _ := lmdatabase.NotificationGetByDepartment(db, department)
	funcMap := template.FuncMap{
		"priorityMap": func(prio int) string {
			priorityMap := map[int]string{
				1: "is-danger",
				2: "is-warning",
				3: "is-info",
			}
			return priorityMap[prio]
		},
		"toTime": func(now int64) string {
			return time.Unix(now, 0).Format("15:04:05")
		},
	}
	data := map[string]interface{}{
		"Notifications": notifications,
	}
	var notificationsBuffer bytes.Buffer

	templateString, err := box.String("templates/notifications.html")
	if err != nil {
		log.Fatal(err)
	}

	notificationHTML := template.Must(template.New("notifications").Funcs(funcMap).Parse(templateString))
	notificationHTML.Execute(&notificationsBuffer, data)
	return notificationsBuffer.String()
}
