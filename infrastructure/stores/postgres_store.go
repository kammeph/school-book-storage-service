package stores

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/kammeph/school-book-storage-service/application/common"
)

const (
	insertSql     = "INSERT INTO ${TABLE} (id, aggregate_id, data, version) VALUES (gen_random_uuid(),$1, $2, $3)"
	selectSql     = "SELECT data, version FROM ${TABLE} WHERE aggregate_id = $1 AND version >= $2 AND version <= $3 ORDER BY version ASC"
	maxVersionSql = "SELECT MAX(version) FROM ${TABLE} WHERE aggregate_id = $1"
)

type PostgresStore struct {
	tableName string
	db        *sql.DB
}

func NewPostgressStore(tableName string, db *sql.DB) common.Store {
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

func (s *PostgresStore) Load(ctx context.Context, aggregateID string) ([]common.Record, error) {
	return s.LoadVersions(ctx, aggregateID, 0, 0)
}

func (s *PostgresStore) LoadVersions(ctx context.Context, aggregateID string, fromVersion int, toVersion int) ([]common.Record, error) {
	stmt, err := s.db.PrepareContext(ctx, s.expand(selectSql))
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	if toVersion == 0 {
		toVersion = math.MaxInt32
	}

	rows, err := stmt.QueryContext(ctx, aggregateID, fromVersion, toVersion)
	if err != nil {
		return nil, err
	}

	records := []common.Record{}
	for rows.Next() {
		record := common.Record{}
		if err := rows.Scan(&record.Data, &record.Version); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *PostgresStore) Save(ctx context.Context, aggregateID string, records ...common.Record) error {
	if len(records) == 0 {
		return nil
	}

	history := common.History(records)
	sort.Sort(history)
	maxVersion, err := s.maxVersion(ctx, aggregateID)
	if err != nil {
		return err
	}

	if maxVersion >= history[0].Version {
		return fmt.Errorf("Event version (%v) is lower than max version (%v) of aggregate", history[0].Version, maxVersion)
	}

	stmt, err := s.db.PrepareContext(ctx, s.expand(insertSql))
	defer stmt.Close()
	if err != nil {
		return err
	}

	for _, record := range history {
		stmt.Exec(aggregateID, record.Data, record.Version)
	}

	return nil
}
