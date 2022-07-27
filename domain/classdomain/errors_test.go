package classdomain_test

import (
	"testing"

	"github.com/kammeph/school-book-storage-service/domain/classdomain"
	"github.com/stretchr/testify/assert"
)

func TestErrApplyEventClassAlreadyExists(t *testing.T) {
	err := classdomain.ErrApplyEventClassAlreadyExists("Test-Type", "3A")
	assert.Error(t, err)
	assert.Equal(t, "can not apply Test-Type: Class with ID 3A already exists", err.Error())
}

func TestErrApplyEventClassNotFound(t *testing.T) {
	err := classdomain.ErrApplyEventClassNotFound("Test-Type", "3A")
	assert.Error(t, err)
	assert.Equal(t, "can not apply Test-Type: Class with ID 3A not found", err.Error())
}

func TestErrClassAlreadyExists(t *testing.T) {
	err := classdomain.ErrClassAlreadyExists(3, "B", 2022, 2023)
	assert.Error(t, err)
	assert.Equal(t, "there is already a similar class 3B (2022-2023)", err.Error())
}

func TestErrClassWithIDNotFound(t *testing.T) {
	err := classdomain.ErrClassWithIDNotFound("5C")
	assert.Error(t, err)
	assert.Equal(t, "class with ID 5C not found", err.Error())
}
