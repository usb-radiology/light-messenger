package lmdatabase

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

// Notification ..
type Notification struct {
	NotificationID string
	DepartmentID   string
	Priority       int
	Modality       string
	CreatedAt      int64
	ConfirmedAt    int64 // default 0, i.e. NULL
	CancelledAt    int64 // default 0, i.e. NULL
}

// NotificationInsert ..
func NotificationInsert(db *sql.DB, department string, priority int, modality string, createdAt int64) error {
	insertStmt, err := db.Prepare(`
	INSERT INTO
		Notification (notificationId, departmentId, priority, modality, createdAt)
	VALUES( ?, ? , ?, ?, ?) `)

	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, errExec := insertStmt.Exec(uuid.New().String(), department, priority, modality, createdAt)
	if errExec != nil {
		return errExec
	}
	return nil
}

// NotificationGetByDepartmentAndModality ..
func NotificationGetByDepartmentAndModality(db *sql.DB, department string, modality string) (*Notification, error) {
	queryStmt :=
		`SELECT
			notificationId, departmentId, modality, priority, createdAt
		FROM
			Notification
		WHERE
			departmentId = ?	
		AND
			modality = ?	
		AND
			cancelledAt = -1
		AND
			confirmedAt = -1`

	row := db.QueryRow(queryStmt, department, modality)
	//defer db.Close()
	var result Notification
	errRowScan := row.Scan(&result.NotificationID, &result.DepartmentID, &result.Modality, &result.Priority, &result.CreatedAt)
	if errRowScan != nil {
		if errRowScan == sql.ErrNoRows {
			result.Modality = modality
			result.DepartmentID = department
			result.Priority = 99
			return &result, nil
		}

		return &result, errRowScan
	}
	return &result, nil
}

// NotificationGetByDepartment ..
func NotificationGetByDepartment(db *sql.DB, department string) (*[]Notification, error) {
	queryStmt :=
		`SELECT
			notificationId, modality, departmentId, priority, createdAt
		FROM
			Notification
		WHERE
			departmentId = ?
		AND
			cancelledAt = -1
		AND
			confirmedAt = -1
		ORDER BY 
			priority 
		ASC`

	rows, errQuery := db.Query(queryStmt, department)
	openNotifications := make([]Notification, 0)
	if errQuery != nil {
		return &openNotifications, errQuery
	}
	defer rows.Close()
	

	for rows.Next() {
		var notification Notification
		if errRowScan := rows.Scan(&notification.NotificationID, &notification.Modality, &notification.DepartmentID, &notification.Priority, &notification.CreatedAt); errRowScan != nil {
			return nil, errRowScan
		}
		openNotifications = append(openNotifications, notification)
	}
	//defer rows.Close()
	return &openNotifications, nil
}

// NotificationGetByModality ..
func NotificationGetByModality(db *sql.DB, modality string) (*[]Notification, error) {
	queryStmt :=
		`SELECT
			notificationId, modality, departmentId, priority, createdAt, confirmedAt, cancelledAt
		FROM
			Notification
		WHERE
			modality = ?
		AND
			(confirmedAt <> -1 OR cancelledAt <> -1)
		ORDER BY 
			createdAt 
		DESC`

	rows, errQuery := db.Query(queryStmt, modality)
	if errQuery != nil {
		log.Fatal(errQuery)
		return nil, errQuery
	}
	defer rows.Close()
	processedNotifications := make([]Notification, 0)

	for rows.Next() {
		var notification Notification
		if errRowScan := rows.Scan(&notification.NotificationID, &notification.Modality,
			&notification.DepartmentID, &notification.Priority, &notification.CreatedAt,
			&notification.ConfirmedAt, &notification.CancelledAt); errRowScan != nil {
			log.Printf("%+v, %+v", notification, processedNotifications)
			return nil, errRowScan
		}
		processedNotifications = append(processedNotifications, notification)
	}
	//defer rows.Close()
	return &processedNotifications, nil
}

// NotificationCancel ..
func NotificationCancel(db *sql.DB, modality string, department string, cancelledAt int64) error {
	updateStmt, err := db.Prepare(`
	UPDATE
		Notification
	SET
		cancelledAt = ?
	WHERE
		modality = ?
	AND
		departmentId = ?
	AND
		confirmedAt = -1
	AND
		cancelledAt = -1`)

	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, errExec := updateStmt.Exec(cancelledAt, modality, department)
	if errExec != nil {
		return errExec
	}

	return nil
}

// NotificationUpdatePriority ..
func NotificationUpdatePriority(db *sql.DB, notificationID string, priority int) error {
	updateStmt, err := db.Prepare(`
	UPDATE
		Notification
	SET
		priority = ?
	WHERE
		notificationId = ?`)

	if err != nil {
		return err
	}
	defer updateStmt.Close()

	_, errExec := updateStmt.Exec(priority, notificationID)
	if errExec != nil {
		return errExec
	}

	return nil
}

// NotificationConfirm ..
func NotificationConfirm(db *sql.DB, notificationID string, now int64) error {
	updateStmt, err := db.Prepare(`
	UPDATE
		Notification
	SET
		confirmedAt = ?
	WHERE
		notificationId = ?`)

	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, errExec := updateStmt.Exec(now, notificationID)
	if errExec != nil {
		return errExec
	}

	return nil
}
