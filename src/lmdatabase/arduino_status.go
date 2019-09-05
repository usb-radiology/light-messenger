package lmdatabase

import "database/sql"

// ArduinoStatus ..
type ArduinoStatus struct {
	DepartmentID string
	StatusAt     int64
}

// ArduinoStatusInsert ..
func ArduinoStatusInsert(db *sql.DB, status ArduinoStatus) error {

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

// ArduinoStatusQueryWithin5MinutesFromNow ..
func ArduinoStatusQueryWithin5MinutesFromNow(db *sql.DB, department string, now int64) (*ArduinoStatus, error) {

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
		return nil, errRowScan
	}

	return &result, nil
}
