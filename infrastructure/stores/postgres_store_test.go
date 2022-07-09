package stores_test

import (
	"context"
	"encoding/json"
	"math"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service/application/common"
	"github.com/kammeph/school-book-storage-service/infrastructure/stores"

	"github.com/stretchr/testify/assert"
)

const (
	insertSql     = "INSERT INTO test \\(aggregate_id, data, version\\) VALUES \\(\\$1, \\$2, \\$3\\)"
	selectSql     = "SELECT data, version FROM test WHERE aggregate_id = \\$1 AND version >= \\$2 AND version <= \\$3 ORDER BY version ASC"
	maxVersionSql = "SELECT MAX\\(version\\) FROM test WHERE aggregate_id = \\$1"
)

func TestNewPostgresStore(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)
	store := stores.NewPostgressStore("test", db)
	assert.NotNil(t, store)
	pgStore, ok := store.(*stores.PostgresStore)
	assert.True(t, ok)
	assert.NotNil(t, pgStore)
}

func TestPostgresStoreLoad(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	store := stores.NewPostgressStore("test", db)

	aggregateID := uuid.New().String()
	rows := sqlmock.NewRows([]string{"data", "version"}).AddRow("My Data", 1).AddRow("Second Data", 2)
	mock.ExpectPrepare(selectSql).ExpectQuery().WithArgs(aggregateID, 0, math.MaxInt32).WillReturnRows(rows)

	ctx := context.Background()
	records, err := store.Load(ctx, aggregateID)
	assert.Nil(t, err)
	assert.Len(t, records, 2)
}

func TestPostgresStoreSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	data, err := json.Marshal("Record1")
	assert.Nil(t, err)
	record := common.Record{Version: 1, Data: data}
	store := stores.NewPostgressStore("test", db)
	aggregateID := uuid.New().String()

	rows := sqlmock.NewRows([]string{"version"})
	mock.ExpectPrepare(maxVersionSql).ExpectQuery().WillReturnRows(rows)

	mock.ExpectPrepare(insertSql).ExpectExec().WithArgs(aggregateID, data, 1)

	ctx := context.Background()
	err = store.Save(ctx, aggregateID, record)
	assert.Nil(t, err)
}

func TestPostgresStoreSaveWithMaxVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	data, err := json.Marshal("Record1")
	assert.Nil(t, err)
	record := common.Record{Version: 4, Data: data}
	store := stores.NewPostgressStore("test", db)
	aggregateID := uuid.New().String()

	rows := sqlmock.NewRows([]string{"version"}).AddRow(3)
	mock.ExpectPrepare(maxVersionSql).ExpectQuery().WillReturnRows(rows)

	mock.ExpectPrepare(insertSql).ExpectExec().WithArgs(aggregateID, data, 1)

	ctx := context.Background()
	err = store.Save(ctx, aggregateID, record)
	assert.Nil(t, err)
}

func TestProstgresSaveWithWrongVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	data, err := json.Marshal("Record")
	assert.Nil(t, err)
	record := common.Record{Version: 1, Data: data}
	store := stores.NewPostgressStore("test", db)
	aggregateID := uuid.New().String()

	rows := sqlmock.NewRows([]string{"version"}).AddRow(2)
	mock.ExpectPrepare(maxVersionSql).ExpectQuery().WillReturnRows(rows)

	mock.ExpectPrepare(insertSql).ExpectExec().WithArgs(aggregateID, data, 1)

	ctx := context.Background()
	err = store.Save(ctx, aggregateID, record)
	assert.NotNil(t, err)
}
