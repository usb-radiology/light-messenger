package lmdatabase

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationShouldInsertAndGetNotification(t *testing.T) {

	// given
	db := setupTest(t)

	departmentID := "abc"
	modality := "def"
	priority := 1
	createdAt := int64(1000)

	// when
	{
		errInsert := NotificationInsert(db, departmentID, priority, modality, createdAt)
		if errInsert != nil {
			t.Fatalf("%+v", errors.WithStack(errInsert))
		}
	}

	// then
	result, errQuery := NotificationGetByDepartmentAndModality(db, departmentID, modality)
	if errQuery != nil {
		t.Fatalf("%+v", errors.WithStack(errQuery))
	}

	if result == nil {
		t.Errorf("Did not retrieve the notification")
	}

	// fmt.Printf("%+v\n", result)
	assert.Equal(t, departmentID, result.DepartmentID)
	assert.Equal(t, modality, result.Modality)
	assert.Equal(t, priority, result.Priority)
	assert.Equal(t, createdAt, result.CreatedAt)
	assert.Equal(t, int64(0), result.ConfirmedAt)

}
