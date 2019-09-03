package lmdatabase

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

// GetDB ..
func GetDB(initConfig *configuration.Configuration) (*sql.DB, error) {
	conn := initConfig.Database.Username + ":" + initConfig.Database.Password + "@tcp(" + initConfig.Database.Host + ":" + strconv.Itoa(initConfig.Database.Port) + ")/" + initConfig.Database.DBName
	log.Print(conn)
	return sql.Open("mysql", conn)
}

// ArduinoStatus ..
type ArduinoStatus struct {
	DepartmentID string
	StatusAt     string
}

// InsertStatus ..
func InsertStatus(db *sql.DB, status ArduinoStatus) error {

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(`
	INSERT INTO 
		ArduinoStatus 
	VALUES( ?, NOW() ) 
		ON DUPLICATE KEY UPDATE 
	statusAt = NOW()`)

	if err != nil {
		return err
	}

	defer stmtIns.Close()
	_, errExec := stmtIns.Exec(status.DepartmentID)
	if errExec != nil {
		return errExec
	}
	return nil
}
