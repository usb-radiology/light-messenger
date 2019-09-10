package server

import (
	"bytes"
	"database/sql"
	"text/template"
	"time"

	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func getCardHTML(db *sql.DB, modality string, department string) (string, error) {
	now := time.Now().Unix()
	arduinoStatus, errStatusQuery := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, now)
	if errStatusQuery != nil {
		return "", errStatusQuery
	}

	notification, errNotificationGetByDepartmentAndModality := lmdatabase.NotificationGetByDepartmentAndModality(db, department, modality)
	if errNotificationGetByDepartmentAndModality != nil {
		return "", errNotificationGetByDepartmentAndModality
	}

	data := map[string]interface{}{
		"Modality":       notification.Modality,
		"Department":     notification.DepartmentID,
		"PriorityNumber": notification.Priority,
		"PriorityName":   priorityMap[notification.Priority],
		"ArduinoStatus":  arduinoStatus,
		"CreatedAt":      time.Unix(notification.CreatedAt, 0).Format("15:04:05"),
	}

	var aodBuffer bytes.Buffer
	errExecute := templates[templateCardID].Execute(&aodBuffer, data)
	if errExecute != nil {
		return "", errExecute
	}

	return aodBuffer.String(), nil
}

func getNotificationsHTML(db *sql.DB, department string) (string, error) {
	notifications, errNotificationGetByDepartment := lmdatabase.NotificationGetByDepartment(db, department)
	if errNotificationGetByDepartment != nil {
		return "", errNotificationGetByDepartment
	}

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
		return "", err
	}

	notificationHTML := template.Must(template.New("notifications").Funcs(funcMap).Parse(templateString))
	errExecute := notificationHTML.Execute(&notificationsBuffer, data)
	if errExecute != nil {
		return "", errExecute
	}

	return notificationsBuffer.String(), nil
}
