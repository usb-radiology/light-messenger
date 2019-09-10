package server

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
	"github.com/usb-radiology/light-messenger/src/version"
)

func mainHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	data := map[string]interface{}{
		"Version":   version.Version,
		"BuildTime": version.BuildTime,
	}

	errRenderTemplate := renderTemplate(w, r, templates[templateIndexID], data)
	if writeInternalServerError(errRenderTemplate, w) != nil {
		return errRenderTemplate
	}

	return nil
}

func visierungHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]

	processedNotifications, errNotificationGetByModality := lmdatabase.NotificationGetByModality(db, modality)
	if writeInternalServerError(errNotificationGetByModality, w) != nil {
		return errNotificationGetByModality
	}

	data := map[string]interface{}{
		"Modality":               modality,
		"AOD":                    getCardHTML(db, modality, "aod"),
		"CTD":                    getCardHTML(db, modality, "ctd"),
		"MSK":                    getCardHTML(db, modality, "msk"),
		"NR":                     getCardHTML(db, modality, "nr"),
		"NUK_NUK":                getCardHTML(db, modality, "nuk"),
		"Version":                version.Version,
		"BuildTime":              version.BuildTime,
		"ProcessedNotifications": processedNotifications,
	}

	errRenderTemplate := renderTemplate(w, r, templates[templateVisierungID], data)
	if writeInternalServerError(errRenderTemplate, w) != nil {
		return errRenderTemplate
	}

	return nil
}

func confirmHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-IC-Remove", "true")

	vars := mux.Vars(r)
	notificationID := vars["id"]

	errNotificationConfirm := lmdatabase.NotificationConfirm(db, notificationID, time.Now().Unix())
	if writeInternalServerError(errNotificationConfirm, w) != nil {
		return errNotificationConfirm
	}

	return nil
}

func radiologieHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	department := vars["department"]

	arduinoStatus, errStatusQuery := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, time.Now().Unix())
	if writeInternalServerError(errStatusQuery, w) != nil {
		return errStatusQuery
	}

	data := map[string]interface{}{
		"Department":    department,
		"Notifications": getNotificationHTML(db, department),
		"Version":       version.Version,
		"BuildTime":     version.BuildTime,
		"ArduinoStatus": arduinoStatus,
	}

	errRenderTemplate := renderTemplate(w, r, templates[templateRadiologieID], data)
	if writeInternalServerError(errRenderTemplate, w) != nil {
		return errRenderTemplate
	}

	return nil
}

func priorityHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]
	priority := vars["priority"]
	log.Print("priorityHandler ", modality, ", ", department, ", ", priority)

	priorityNumber, errPriorityConversion := strconv.Atoi(priority)
	if writeInternalServerError(errPriorityConversion, w) != nil {
		return errPriorityConversion
	}

	notification, errNotificationGetByDepartmentAndModality := lmdatabase.NotificationGetByDepartmentAndModality(db, department, modality)
	if writeInternalServerError(errNotificationGetByDepartmentAndModality, w) != nil {
		return errNotificationGetByDepartmentAndModality
	}

	now := time.Now().Unix()

	if notification.NotificationID == "" {
		errNotificationInsert := lmdatabase.NotificationInsert(db, department, priorityNumber, modality, now)
		if writeInternalServerError(errNotificationInsert, w) != nil {
			return errNotificationInsert
		}

	} else {

		errNotificationUpdatePriority := lmdatabase.NotificationUpdatePriority(db, notification.NotificationID, priorityNumber)
		if writeInternalServerError(errNotificationUpdatePriority, w) != nil {
			return errNotificationUpdatePriority
		}

	}

	arduinoStatus, errStatusQuery := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, now)
	if writeInternalServerError(errStatusQuery, w) != nil {
		return errStatusQuery
	}

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"Priority":       priority,
		"PriorityName":   priorityMap[priorityNumber],
		"PriorityNumber": priorityNumber,
		"ArduinoStatus":  arduinoStatus,
		"CreatedAt":      time.Unix(now, 0).Format("15:04:05"),
	}

	errRenderTemplateName := renderTemplateName(w, r, templates[templateCardID], "card_view", data)
	if writeInternalServerError(errRenderTemplateName, w) != nil {
		return errRenderTemplateName
	}

	return nil
}

func cancelHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"PriorityNumber": 99, // needed because of le comparison in template
	}

	errNotificationCancel := lmdatabase.NotificationCancel(db, modality, department, time.Now().Unix())
	if writeInternalServerError(errNotificationCancel, w) != nil {
		return errNotificationCancel
	}

	errRenderTemplateName := renderTemplateName(w, r, templates[templateCardID], "card_view", data)
	if writeInternalServerError(errRenderTemplateName, w) != nil {
		return errRenderTemplateName
	}

	return nil
}
