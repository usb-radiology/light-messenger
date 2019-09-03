package lmdatabase

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // mysql driver ..
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

// GetDB ..
func GetDB(initConfig *configuration.Configuration) (*sql.DB, error) {
	conn := initConfig.Database.Username + ":" + initConfig.Database.Password + "@tcp(" + initConfig.Database.Host + ":" + strconv.Itoa(initConfig.Database.Port) + ")/" + initConfig.Database.DBName
	return sql.Open("mysql", conn)
}

func execStatements(db *sql.DB, sqlStatements []string) error {

	for _, statement := range sqlStatements {
		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

// ArduinoStatus ..
type ArduinoStatus struct {
	DepartmentID string
	StatusAt     int64
}

// Notification ..
type Notification struct {
	NotificationID string
	DepartmentID   string
	Priority       int
	Modality       string
	CreatedAt      int64
	ConfirmedAt    int64
}

// InsertStatus ..
func InsertStatus(db *sql.DB, status ArduinoStatus) error {

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(`
	INSERT INTO 
		ArduinoStatus 
	VALUES( ?, ? ) 
		ON DUPLICATE KEY UPDATE 
	statusAt =?`)

	if err != nil {
		return err
	}

	defer stmtIns.Close()
	_, errExec := stmtIns.Exec(status.DepartmentID, status.StatusAt, status.StatusAt)
	if errExec != nil {
		return errExec
	}
	return nil
}

// IsAlive ..
func IsAlive(db *sql.DB, department string, now int64) (*ArduinoStatus, error) {

	queryStr := `
	SELECT departmentId, statusAt FROM 
		ArduinoStatus 
	WHERE
		departmentId = ?
	AND 
		statusAt > ?`

	row := db.QueryRow(queryStr, department, now-300)

	var result ArduinoStatus

	errRowScan := row.Scan(&result.DepartmentID, &result.StatusAt)
	if errRowScan != nil {
		if errRowScan == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(errRowScan, "error retrieving arduino status")
	}

	return &result, nil
}

// InsertNotification ..
func InsertNotification(db *sql.DB, department string, priority int, modality string, createdAt int64) error {
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(`
	INSERT INTO 
		Notification (notificationId, departmentId, priority, modality, createdAt)
	VALUES( ?, ? , ?, ?, ?) `)

	if err != nil {
		return err
	}

	defer stmtIns.Close()
	_, errExec := stmtIns.Exec(uuid.New().String(), department, priority, modality, createdAt)
	if errExec != nil {
		return errors.Wrap(errExec, "error inserting notification")
	}
	return nil
}

// SelectNotification ..
func SelectNotification(db *sql.DB, department string, priority int, modality string) (*Notification, error) {
	// Prepare statement for inserting data
	queryStr :=
		`SELECT 
			notificationId, departmentId, priority, modality, createdAt
		FROM 
			Notification 
		WHERE
			confirmedAt IS NULL`

	row := db.QueryRow(queryStr)

	var result Notification

	errRowScan := row.Scan(&result.NotificationID, &result.DepartmentID, &result.DepartmentID, &result.Priority, &result.Modality, &result.CreatedAt)
	if errRowScan != nil {
		if errRowScan == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(errRowScan, "error retrieving notification")
	}

	return &result, nil
}


// UpdateNotification ..
func UpdateNotification(db *sql.DB, notificationID string, priority int) error {
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(`
	UPDATE
		Notification 
	SET
		priority = ?
	WHERE 
		notificationId = ?`)

	if err != nil {
		return err
	}

	defer stmtIns.Close()
	_, errExec := stmtIns.Exec(priority, notificationID)
	if errExec != nil {
		return errors.Wrap(errExec, "error updating notification")
	}
	return nil
}