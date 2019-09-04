package server

import (
	"bytes"
	"database/sql"
	"text/template"
	"time"

	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func create(db *sql.DB, modality string, department string) string {

	notification, _ := lmdatabase.QueryNotification(db, modality, department)

	// AOD Card Template
	data := map[string]interface{}{
		"Modality":       notification.Modality,
		"Department":     notification.DepartmentID,
		"PriorityNumber": notification.Priority,
		"PriorityName":   priorityName(notification.Priority),
		"CreatedAt":      time.Unix(notification.CreatedAt, 0).Format("15:04:05"),
	}
	var aodBuffer bytes.Buffer
	aodCard := template.Must(template.ParseFiles("templates/card.html"))
	aodCard.Execute(&aodBuffer, data)
	return aodBuffer.String()
}

func priorityName(priority int) string {
	priorityMap := map[int]string{
		1: "is-danger",
		2: "is-warning",
		3: "is-info",
	}
	return priorityMap[priority]
}
