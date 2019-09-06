package lmdatabase

import (
	"database/sql"

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
			cancelledAt IS NULL
		AND
			confirmedAt IS NULL`

	row := db.QueryRow(queryStmt, department, modality)

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
			cancelledAt IS NULL
		AND
			confirmedAt IS NULL`

	rows, errQuery := db.Query(queryStmt, department)
	if errQuery != nil {
		return nil, errQuery
	}

	openNotifications := make([]Notification, 0)

	for rows.Next() {
		var notification Notification
		if errRowScan := rows.Scan(&notification.NotificationID, &notification.Modality, &notification.DepartmentID, &notification.Priority, &notification.CreatedAt); errRowScan != nil {
			return nil, errRowScan
		}
		openNotifications = append(openNotifications, notification)
	}

	return &openNotifications, nil
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
		confirmedAt IS NULL
	AND
		cancelledAt IS NULL`)

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
