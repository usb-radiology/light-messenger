package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func arduinoStatusHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	department := vars["department"]

	status := lmdatabase.ArduinoStatus{
		DepartmentID: department,
		StatusAt:     time.Now().Unix(),
	}

	{
		errInsert := lmdatabase.ArduinoStatusInsert(db, status)
		if writeInternalServerError(errInsert, w) != nil {
			return errInsert
		}
	}

	{
		errWrite := writeBytes(w, []byte(fmt.Sprintf("%+v", status)))
		if writeInternalServerError(errWrite, w) != nil {
			return errWrite
		}
	}

	return nil
}

func openStatusHandler(config *configuration.Configuration, db *sql.DB, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	department := vars["department"]

	notifications, err := lmdatabase.NotificationGetByDepartment(db, department)
	if writeInternalServerError(err, w) != nil {
		return err
	}

	if len(*notifications) > 0 {
		arduinoPrioMap := map[int]string{
			1: "HIGH",
			2: "MEDIUM",
			3: "LOW",
		}

		{
			errWrite := writeBytes(w, []byte(fmt.Sprintf(";1;%v;", arduinoPrioMap[(*notifications)[0].Priority])))
			if writeInternalServerError(errWrite, w) != nil {
				return errWrite
			}
		}

	} else {

		{
			errWrite := writeBytes(w, []byte(fmt.Sprintf(";0;")))
			if writeInternalServerError(errWrite, w) != nil {
				return errWrite
			}
		}
	}

	return nil
}
