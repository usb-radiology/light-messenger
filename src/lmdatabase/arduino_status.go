package lmdatabase

import (
	"database/sql"

	"github.com/pkg/errors"
)

// ArduinoStatus ..
type ArduinoStatus struct {
	DepartmentID string
	StatusAt     int64
}

// ArduinoStatusInsert ..
func ArduinoStatusInsert(db *sql.DB, status ArduinoStatus) error {

	insertStmt, err := db.Prepare(`
	INSERT INTO
		ArduinoStatus
	VALUES( ?, ? )
		ON DUPLICATE KEY UPDATE
	statusAt =?`)

	if err != nil {
		return errors.WithStack(err)
	}

	defer insertStmt.Close()

	_, errExec := insertStmt.Exec(status.DepartmentID, status.StatusAt, status.StatusAt)
	if errExec != nil {
		return errors.WithStack(errExec)
	}

	return nil
}

// ArduinoStatusQueryWithin5MinutesFromNow ..
func ArduinoStatusQueryWithin5MinutesFromNow(db *sql.DB, department string, now int64) (*ArduinoStatus, error) {

	queryStmt := `
	SELECT
		departmentId, statusAt
	FROM
		ArduinoStatus
	WHERE
		departmentId = ?
	AND
		statusAt > ?`

	row := db.QueryRow(queryStmt, department, now-300)

	var result ArduinoStatus

	errRowScan := row.Scan(&result.DepartmentID, &result.StatusAt)
	if errRowScan != nil {
		if errRowScan == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.WithStack(errRowScan)
	}
	return &result, nil
}
