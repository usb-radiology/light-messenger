package lmdatabase

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

// GetDB ..
func GetDB(initConfig *configuration.Configuration) (*sql.DB, error) {
	conn := initConfig.Database.Username + ":" + initConfig.Database.Password + "@tcp(" + initConfig.Database.Host + ":" + strconv.Itoa(initConfig.Database.Port) + ")/" + initConfig.Database.DBName
	return sql.Open("mysql", conn)
}

// ArduinoStatus ..
type ArduinoStatus struct {
	DepartmentID string
	StatusAt     int64
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

// IsAslive ...
func IsAlive(config *configuration.Configuration, department string, now int64) error {
	// Prepare statement for inserting data
	db, _ := GetDB(config)

	stmtIns, err := db.Prepare(`
	SELECT * FROM 
		ArduinoStatus 
	WHERE
		departmentId = ?
	AND 
		statusAt > ?`)

	if err != nil {
		return err
	}

	defer stmtIns.Close()
	_, errExec := stmtIns.Exec(department, now - 300)
	if errExec != nil {
		return errExec
	}
	return nil
}
