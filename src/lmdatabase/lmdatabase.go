package lmdatabase

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // mysql driver ..
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
		return nil, errors.Wrap(errRowScan, "error retrieving result")
	}

	return &result, nil
}
