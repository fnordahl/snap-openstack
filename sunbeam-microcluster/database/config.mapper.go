package database

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/canonical/lxd/lxd/db/query"
	"github.com/canonical/lxd/shared/api"
	"github.com/canonical/microcluster/cluster"
)

var _ = api.ServerEnvironment{}

var configItemObjects = cluster.RegisterStmt(`
SELECT config.id, config.key, config.value
  FROM config
  ORDER BY config.key
`)

var configItemObjectsByKey = cluster.RegisterStmt(`
SELECT config.id, config.key, config.value
  FROM config
  WHERE ( config.key = ? )
  ORDER BY config.key
`)

var configItemID = cluster.RegisterStmt(`
SELECT config.id FROM config
  WHERE config.key = ?
`)

var configItemCreate = cluster.RegisterStmt(`
INSERT INTO config (key, value)
  VALUES (?, ?)
`)

var configItemDeleteByKey = cluster.RegisterStmt(`
DELETE FROM config WHERE key = ?
`)

var configItemUpdate = cluster.RegisterStmt(`
UPDATE config
  SET key = ?, value = ?
 WHERE id = ?
`)

// configItemColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the ConfigItem entity.
func configItemColumns() string {
	return "config.id, config.key, config.value"
}

// getConfigItems can be used to run handwritten sql.Stmts to return a slice of objects.
func getConfigItems(ctx context.Context, stmt *sql.Stmt, args ...any) ([]ConfigItem, error) {
	objects := make([]ConfigItem, 0)

	dest := func(scan func(dest ...any) error) error {
		c := ConfigItem{}
		err := scan(&c.ID, &c.Key, &c.Value)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"config\" table: %w", err)
	}

	return objects, nil
}

// getConfigItems can be used to run handwritten query strings to return a slice of objects.
func getConfigItemsRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]ConfigItem, error) {
	objects := make([]ConfigItem, 0)

	dest := func(scan func(dest ...any) error) error {
		c := ConfigItem{}
		err := scan(&c.ID, &c.Key, &c.Value)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"config\" table: %w", err)
	}

	return objects, nil
}

// GetConfigItems returns all available ConfigItems.
// generator: ConfigItem GetMany
func GetConfigItems(ctx context.Context, tx *sql.Tx, filters ...ConfigItemFilter) ([]ConfigItem, error) {
	var err error

	// Result slice.
	objects := make([]ConfigItem, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = cluster.Stmt(tx, configItemObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"configItemObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Key != nil {
			args = append(args, []any{filter.Key}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, configItemObjectsByKey)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"configItemObjectsByKey\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(configItemObjectsByKey)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"configItemObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Key == nil {
			return nil, fmt.Errorf("Cannot filter on empty ConfigItemFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getConfigItems(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getConfigItemsRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"config\" table: %w", err)
	}

	return objects, nil
}

// GetConfigItem returns the ConfigItem with the given key.
// generator: ConfigItem GetOne
func GetConfigItem(ctx context.Context, tx *sql.Tx, key string) (*ConfigItem, error) {
	filter := ConfigItemFilter{}
	filter.Key = &key

	objects, err := GetConfigItems(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"config\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "ConfigItem not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"config\" entry matches")
	}
}

// GetConfigItemID return the ID of the ConfigItem with the given key.
// generator: ConfigItem ID
func GetConfigItemID(ctx context.Context, tx *sql.Tx, key string) (int64, error) {
	stmt, err := cluster.Stmt(tx, configItemID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"configItemID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, key)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "ConfigItem not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"config\" ID: %w", err)
	}

	return id, nil
}

// ConfigItemExists checks if a ConfigItem with the given key exists.
// generator: ConfigItem Exists
func ConfigItemExists(ctx context.Context, tx *sql.Tx, key string) (bool, error) {
	_, err := GetConfigItemID(ctx, tx, key)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateConfigItem adds a new ConfigItem to the database.
// generator: ConfigItem Create
func CreateConfigItem(ctx context.Context, tx *sql.Tx, object ConfigItem) (int64, error) {
	// Check if a ConfigItem with the same key exists.
	exists, err := ConfigItemExists(ctx, tx, object.Key)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"config\" entry already exists")
	}

	args := make([]any, 2)

	// Populate the statement arguments.
	args[0] = object.Key
	args[1] = object.Value

	// Prepared statement to use.
	stmt, err := cluster.Stmt(tx, configItemCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"configItemCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"config\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"config\" entry ID: %w", err)
	}

	return id, nil
}

// DeleteConfigItem deletes the ConfigItem matching the given key parameters.
// generator: ConfigItem DeleteOne-by-Key
func DeleteConfigItem(_ context.Context, tx *sql.Tx, key string) error {
	stmt, err := cluster.Stmt(tx, configItemDeleteByKey)
	if err != nil {
		return fmt.Errorf("Failed to get \"configItemDeleteByKey\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(key)
	if err != nil {
		return fmt.Errorf("Delete \"config\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "ConfigItem not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d ConfigItem rows instead of 1", n)
	}

	return nil
}

// UpdateConfigItem updates the ConfigItem matching the given key parameters.
// generator: ConfigItem Update
func UpdateConfigItem(ctx context.Context, tx *sql.Tx, key string, object ConfigItem) error {
	id, err := GetConfigItemID(ctx, tx, key)
	if err != nil {
		return err
	}

	stmt, err := cluster.Stmt(tx, configItemUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"configItemUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Key, object.Value, id)
	if err != nil {
		return fmt.Errorf("Update \"config\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}
