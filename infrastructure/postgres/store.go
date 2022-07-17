package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"

	application "github.com/kammeph/school-book-storage-service/application/common"
	domain "github.com/kammeph/school-book-storage-service/domain/common"
)

const (
	insertSql     = "INSERT INTO ${TABLE} (id, aggregate_id, type, version, timestamp, data) VALUES (gen_random_uuid(),$1, $2, $3, $4, $5)"
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
	defer stmt.Close()
	if err != nil {
		return -1, err
	}

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

func (s *PostgresStore) Load(ctx context.Context, aggregate domain.Aggregate) error {
	return s.loadVersions(ctx, aggregate, 0, 0)
}

func (s *PostgresStore) loadVersions(ctx context.Context, aggregate domain.Aggregate, fromVersion int, toVersion int) error {
	stmt, err := s.db.PrepareContext(ctx, s.expand(selectSql))
	defer stmt.Close()
	if err != nil {
		return err
	}

	if toVersion == 0 {
		toVersion = math.MaxInt32
	}

	rows, err := stmt.QueryContext(ctx, aggregate.AggregateID(), fromVersion, toVersion)
	if err != nil {
		return err
	}

	var events []domain.Event
	for rows.Next() {
		event := domain.EventModel{}
		if err := rows.Scan(&event.ID, &event.Type, &event.Version, &event.At, &event.Data); err != nil {
			return err
		}
		events = append(events, &event)
	}

	if err := aggregate.Load(events); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) Save(ctx context.Context, aggregate domain.Aggregate) error {
	if len(aggregate.DomainEvents()) == 0 {
		return nil
	}

	history := domain.History(aggregate.DomainEvents())
	sort.Sort(history)
	maxVersion, err := s.maxVersion(ctx, aggregate.AggregateID())
	if err != nil {
		return err
	}

	if maxVersion >= history[0].EventVersion() {
		return fmt.Errorf("Event version (%v) is lower than max version (%v) of aggregate", history[0].EventVersion(), maxVersion)
	}

	stmt, err := s.db.PrepareContext(ctx, s.expand(insertSql))
	defer stmt.Close()
	if err != nil {
		return err
	}

	for _, event := range history {
		_, err = stmt.Exec(event.AggregateID(), event.EventType(), event.EventVersion(), event.EventAt(), event.EventData())
		if err != nil {
			return err
		}
	}

	return nil
}
