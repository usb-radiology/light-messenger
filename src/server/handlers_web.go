package server

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
	"github.com/usb-radiology/light-messenger/src/version"
)

//
// index
//

func mainHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	data := map[string]interface{}{
		"Version":   version.Version,
		"BuildTime": version.BuildTime,
	}

	return renderTemplate(w, r, templates[templateIndexID], data)
}

//
// MTRA
//

func visierungHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]

	processedNotifications, errNotificationGetByModality := lmdatabase.NotificationGetProcessedNotificationsByModality(db, modality)
	if errNotificationGetByModality != nil {
		return errNotificationGetByModality
	}

	aodCardHTML, errAodCardHTML := getCardHTML(db, modality, "aod")
	if errAodCardHTML != nil {
		return errAodCardHTML
	}

	ctdCardHTML, errCtdCardHTML := getCardHTML(db, modality, "ctd")
	if errCtdCardHTML != nil {
		return errCtdCardHTML
	}

	mskCardHTML, errMskCardHTML := getCardHTML(db, modality, "msk")
	if errMskCardHTML != nil {
		return errMskCardHTML
	}

	nrCardHTML, errNrCardHTML := getCardHTML(db, modality, "nr")
	if errNrCardHTML != nil {
		return errNrCardHTML
	}

	nukCardHTML, errNukCardHTML := getCardHTML(db, modality, "nuk")
	if errNukCardHTML != nil {
		return errNukCardHTML
	}

	data := map[string]interface{}{
		"Modality":               modality,
		"AOD":                    aodCardHTML,
		"CTD":                    ctdCardHTML,
		"MSK":                    mskCardHTML,
		"NR":                     nrCardHTML,
		"NUK_NUK":                nukCardHTML,
		"Version":                version.Version,
		"BuildTime":              version.BuildTime,
		"ProcessedNotifications": processedNotifications,
	}

	if r.Header.Get(HTMLHeaderContentType) == HTMLHeaderContentTypeValueJSON {
		return writeJSON(w, data)
	}

	return renderTemplate(w, r, templates[templateVisierungID], data)
}

//
// Radiology
//

func radiologieHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	department := vars["department"]

	arduinoStatus, errStatusQuery := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, time.Now().Unix())
	if errStatusQuery != nil {
		return errStatusQuery
	}

	notificationsHTML, errNotificationsHTML := getNotificationsHTML(db, department)
	if errNotificationsHTML != nil {
		return errNotificationsHTML
	}

	data := map[string]interface{}{
		"Department":    department,
		"Notifications": notificationsHTML,
		"Version":       version.Version,
		"BuildTime":     version.BuildTime,
		"ArduinoStatus": arduinoStatus,
	}

	if r.Header.Get(HTMLHeaderContentType) == HTMLHeaderContentTypeValueJSON {
		return writeJSON(w, data)
	}

	return renderTemplate(w, r, templates[templateRadiologieID], data)
}

//
// Notifications
//

func notificationCreateHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]
	priority := vars["priority"]
	log.Print("priorityHandler ", modality, ", ", department, ", ", priority)

	priorityNumber, errPriorityConversion := strconv.Atoi(priority)
	if errPriorityConversion != nil {
		return errors.WithStack(errPriorityConversion)
	}

	notification, errNotificationGetByDepartmentAndModality := lmdatabase.NotificationGetOpenNotificationByDepartmentAndModality(db, department, modality)
	if errNotificationGetByDepartmentAndModality != nil {
		return errNotificationGetByDepartmentAndModality
	}

	now := time.Now().Unix()

	if notification.NotificationID == "" {
		errNotificationInsert := lmdatabase.NotificationInsert(db, department, priorityNumber, modality, now)
		if errNotificationInsert != nil {
			return errNotificationInsert
		}

	} else {

		errNotificationUpdatePriority := lmdatabase.NotificationUpdatePriority(db, notification.NotificationID, priorityNumber)
		if errNotificationUpdatePriority != nil {
			return errNotificationUpdatePriority
		}

	}

	arduinoStatus, errStatusQuery := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, department, now)
	if errStatusQuery != nil {
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

	if r.Header.Get(HTMLHeaderContentType) == HTMLHeaderContentTypeValueJSON {
		return writeJSON(w, data)
	}

	return renderTemplateName(w, r, templates[templateCardID], "card_view", data)
}

func notificationCancelHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	modality := vars["modality"]
	department := vars["department"]

	data := map[string]interface{}{
		"Modality":       modality,
		"Department":     department,
		"PriorityNumber": 99, // needed because of le comparison in template
	}

	errNotificationCancel := lmdatabase.NotificationCancel(db, modality, department, time.Now().Unix())
	if errNotificationCancel != nil {
		return errNotificationCancel
	}

	if r.Header.Get(HTMLHeaderContentType) == HTMLHeaderContentTypeValueJSON {
		return writeJSON(w, data)
	}

	return renderTemplateName(w, r, templates[templateCardID], "card_view", data)
}

func notificationConfirmHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-IC-Remove", "true")

	vars := mux.Vars(r)
	notificationID := vars["id"]

	rowsAffected, errNotificationConfirm := lmdatabase.NotificationConfirm(db, notificationID, time.Now().Unix())
	if errNotificationConfirm != nil {
		return errNotificationConfirm
	}

	if rowsAffected == 0 {
		writeBadRequest(w)
	}

	return nil
}
