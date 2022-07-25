package postgresdb

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/kammeph/school-book-storage-service/application"
	"github.com/kammeph/school-book-storage-service/domain"
)

const (
	insertSql     = "INSERT INTO ${TABLE} (id, aggregate_id, type, version, timestamp, data) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)"
	selectSql     = "SELECT aggregate_id, type, version, timestamp, data FROM ${TABLE} WHERE aggregate_id = $1 AND version >= $2 AND version <= $3 ORDER BY version ASC"
	maxVersionSql = "SELECT MAX(version) FROM ${TABLE} WHERE aggregate_id = $1"
)

type PostgresStore struct {
	tableName string
	db        *sql.DB
}

func NewPostgressStore(tableName string, db *sql.DB) application.Store {
	return &PostgresStore{tableName: tableName, db: db}
}

func (s *PostgresStore) expand(stmt string) string {
	return strings.Replace(stmt, "${TABLE}", s.tableName, -1)
}

func (s *PostgresStore) maxVersion(ctx context.Context, aggregateID string) (int, error) {
	stmt, err := s.db.PrepareContext(ctx, s.expand(maxVersionSql))
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, aggregateID)
	if err != nil {
		return -1, err
	}

	maxVersion := -1
	for rows.Next() {
		if err := rows.Scan(&maxVersion); err != nil {
			return -1, nil
		}
	}

	return maxVersion, nil
}

func (s *PostgresStore) Load(ctx context.Context, aggregateID string) ([]domain.Event, error) {
	return s.loadVersions(ctx, aggregateID, 0, 0)
}

func (s *PostgresStore) loadVersions(ctx context.Context, aggregateID string, fromVersion int, toVersion int) ([]domain.Event, error) {
	stmt, err := s.db.PrepareContext(ctx, s.expand(selectSql))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if toVersion == 0 {
		toVersion = math.MaxInt32
	}

	rows, err := stmt.QueryContext(ctx, aggregateID, fromVersion, toVersion)
	if err != nil {
		return nil, err
	}

	var events []domain.Event
	for rows.Next() {
		event := domain.EventModel{}
		if err := rows.Scan(&event.ID, &event.Type, &event.Version, &event.At, &event.Data); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (s *PostgresStore) Save(ctx context.Context, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	history := domain.History(events)
	sort.Sort(history)
	maxVersion, err := s.maxVersion(ctx, events[0].AggregateID())
	if err != nil {
		return err
	}

	if maxVersion >= history[0].EventVersion() {
		return fmt.Errorf("event version (%v) is lower than max version (%v) of aggregate", history[0].EventVersion(), maxVersion)
	}

	stmt, err := s.db.PrepareContext(ctx, s.expand(insertSql))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range history {
		_, err = stmt.ExecContext(ctx, event.AggregateID(), event.EventType(), event.EventVersion(), event.EventAt(), event.EventData())
		if err != nil {
			return err
		}
	}

	return nil
}
