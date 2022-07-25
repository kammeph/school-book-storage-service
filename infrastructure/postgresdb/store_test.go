package postgresdb_test

import (
	"context"
	"database/sql/driver"
	"math"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kammeph/school-book-storage-service/domain"
	"github.com/kammeph/school-book-storage-service/infrastructure/postgresdb"
	"github.com/stretchr/testify/assert"
)

const (
	insertSql     = "INSERT INTO test \\(id, aggregate_id, type, version, timestamp, data\\) VALUES \\(gen_random_uuid\\(\\), \\$1, \\$2, \\$3, \\$4, \\$5\\)"
	selectSql     = "SELECT aggregate_id, type, version, timestamp, data FROM test WHERE aggregate_id = \\$1 AND version >= \\$2 AND version <= \\$3 ORDER BY version ASC"
	maxVersionSql = "SELECT MAX\\(version\\) FROM test WHERE aggregate_id = \\$1"
)

func TestNewPostgresStore(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)
	store := postgresdb.NewPostgresStore("test", db)
	assert.NotNil(t, store)
	pgStore, ok := store.(*postgresdb.PostgresStore)
	assert.True(t, ok)
	assert.NotNil(t, pgStore)
}

func TestLoad(t *testing.T) {
	db, mock, _ := sqlmock.New()
	store := postgresdb.NewPostgresStore("test", db)

	aggregateID := "testSchool"
	rows := sqlmock.
		NewRows([]string{"aggregate_id", "type", "version", "timestamp", "data"}).
		AddRow(aggregateID, "testType", 1, time.Now(), "my first data").
		AddRow(aggregateID, "testType", 2, time.Now(), "my second data")
	mock.ExpectPrepare(selectSql).ExpectQuery().WithArgs(aggregateID, 0, math.MaxInt32).WillReturnRows(rows)

	events, err := store.Load(context.Background(), aggregateID)
	assert.Nil(t, err)
	assert.Len(t, events, 2)
}

func TestSave(t *testing.T) {
	tests := []struct {
		name          string
		latestVersion int
		eventVerion   int
		expectError   bool
	}{
		{
			name:          "first save",
			latestVersion: 0,
			eventVerion:   1,
			expectError:   false,
		},
		{
			name:          "correct version order",
			latestVersion: 5,
			eventVerion:   6,
			expectError:   false,
		},
		{
			name:          "incorrect version order",
			latestVersion: 9,
			eventVerion:   5,
			expectError:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			store := postgresdb.NewPostgresStore("test", db)
			rows := sqlmock.NewRows([]string{"version"}).AddRow(test.latestVersion)
			mock.ExpectPrepare(maxVersionSql).ExpectQuery().WillReturnRows(rows)

			event := domain.EventModel{
				ID:      "testSchool",
				Type:    "testType",
				Version: test.eventVerion,
				At:      time.Now(),
				Data:    "my data",
			}
			mock.
				ExpectPrepare(insertSql).
				ExpectExec().
				WithArgs(event.AggregateID(), event.EventType(), event.EventVersion(), event.EventAt(), event.EventData()).
				WillReturnResult(driver.RowsAffected(1))

			err := store.Save(context.Background(), []domain.Event{&event})
			if test.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}

}
