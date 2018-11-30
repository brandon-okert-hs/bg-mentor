package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// A Database is the only public interface to the database
// It internally maintains a reference to the actual database connection, and manages concurrency etc.
// Only one Database needs to exist - multiple databases are redundant as a single Database represents multiple connections
type Database struct {
	config mysql.Config
	sql    *sql.DB
}

// NewDatabase creates a Database with the standard config
func NewDatabase(driver string, config mysql.Config) (*Database, error) {
	db, err := sql.Open(driver, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("Failed to create Database: %s", err)
	}

	return &Database{
		config: config,
		sql:    db,
	}, nil
}

// Amiup returns true if the database is up and ready for connections, false otherwise
func (db *Database) Amiup() bool {
	err := db.sql.Ping()
	if err != nil {
		return false
	}
	return true
}

// NumOpenConnections reports the number of open connections. Useful for health checking and debugging
func (db *Database) NumOpenConnections() int {
	return db.sql.Stats().OpenConnections
}

// SetMaxIdleConns determines how many connects to keep alive for pending requests, even after closing them
// Defaults to 2 internally
func (db *Database) SetMaxIdleConns(n int) {
	db.sql.SetMaxIdleConns(n)
}

// Execute runs an arbitrary query that doesn't return rows with the given params injected in safely
func (db *Database) Execute(query string, params ...interface{}) (sql.Result, error) {
	logger.Debugw("Attempting to execute non-row-query", "query", query, "params", params)

	res, err := db.sql.Exec(query, params...)
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %s", err.Error())
	}

	return res, err
}

// Query runs an arbitrary query that returns rows with the given params injected in safely
func (db *Database) Query(query string, params ...interface{}) (*sql.Rows, error) {
	logger.Debugw("Attempting to execute row-query", "query", query, "params", params)

	rows, err := db.sql.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %s", err.Error())
	}

	return rows, err
}

func nullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func nullInt(i int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(i),
		Valid: i != 0,
	}
}
