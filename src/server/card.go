package server

import (
	"bytes"
	"database/sql"
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
